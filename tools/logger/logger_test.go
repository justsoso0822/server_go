package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"app/config"
)

// TestReplaceAttrTimeLayout 锁定日志时间格式：秒之后再补 3 位毫秒。
// 没有毫秒时，同一秒内的多条日志在 JSON 里 time 完全相同，只能靠行序区分。
func TestReplaceAttrTimeLayout(t *testing.T) {
	// 纳秒部分取 123456789，验证是截断到毫秒而不是四舍五入。
	now := time.Date(2026, 11, 25, 18, 3, 51, 123456789, time.Local)
	got := replaceAttr(nil, slog.Time(slog.TimeKey, now))

	if got.Key != slog.TimeKey {
		t.Fatalf("字段名被改掉了: %q", got.Key)
	}
	const want = "2026-11-25 18:03:51.123"
	if text := got.Value.String(); text != want {
		t.Fatalf("时间格式不对:\n got %q\nwant %q", text, want)
	}
	if _, err := time.Parse(logTimeLayout, got.Value.String()); err != nil {
		t.Fatalf("无法按 logTimeLayout 解析回去: %v", err)
	}
}

// TestNewTruncatesEveryFile 覆盖“多个输出文件 + clearOnStart”场景：
// 每个文件都必须被截断，而不是只有第一个。
func TestNewTruncatesEveryFile(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.log")
	p2 := filepath.Join(dir, "b.log")
	for _, p := range []string{p1, p2} {
		if err := os.WriteFile(p, []byte("old-content\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	log, closeFn, err := New(config.LoggerConfig{
		Level:        "debug",
		Format:       "json",
		ClearOnStart: true,
		OutputPaths:  []string{p1, p2},
	})
	if err != nil {
		t.Fatalf("New 返回错误: %v", err)
	}
	log.Info("fresh")
	if err := closeFn(); err != nil {
		t.Fatalf("closeFn 返回错误: %v", err)
	}

	for _, p := range []string{p1, p2} {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		got := string(b)
		if strings.Contains(got, "old-content") {
			t.Fatalf("%s 没有被清空，内容仍含旧日志: %q", filepath.Base(p), got)
		}
		if !strings.Contains(got, "fresh") {
			t.Fatalf("%s 没有写入新日志: %q", filepath.Base(p), got)
		}
	}
}

// TestOpenWritersDefaultsToStdout 覆盖“一项都没配”的兜底。
// 删掉 Stdout 开关后，stdout 只能来自 outputPaths 或这个兜底。
func TestOpenWritersDefaultsToStdout(t *testing.T) {
	writers, closers, err := openWriters(config.LoggerConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if len(closers) != 0 {
		t.Fatalf("stdout 不应产生需要关闭的文件，得到 %d 个", len(closers))
	}
	if len(writers) != 1 {
		t.Fatalf("空配置应只兜底 1 个 writer，得到 %d 个", len(writers))
	}
	if writers[0] != os.Stdout {
		t.Fatal("空配置兜底的不是 os.Stdout")
	}
}

// TestOpenWritersSpecialValuesIgnoreCase 覆盖特殊值的大小写：
// "STDOUT" 必须写标准输出，而不是在工作目录里建一个叫 STDOUT 的文件。
func TestOpenWritersSpecialValuesIgnoreCase(t *testing.T) {
	t.Chdir(t.TempDir()) // 万一实现写错，别在包目录里留下垃圾文件

	for _, name := range []string{"stdout", "STDOUT", "StdOut"} {
		writers, closers, err := openWriters(config.LoggerConfig{OutputPaths: []string{name}})
		if err != nil {
			t.Fatalf("%q 应被当作特殊值: %v", name, err)
		}
		if len(closers) != 0 {
			closeFiles(closers) // 失败路径也别留占用，Windows 上会挡住 TempDir 清理
			t.Fatalf("%q 不应打开文件", name)
		}
		if len(writers) != 1 || writers[0] != os.Stdout {
			t.Fatalf("%q 应映射到 os.Stdout", name)
		}
	}
	if _, err := os.Stat("STDOUT"); err == nil {
		t.Fatal("不应在当前目录创建 STDOUT 文件")
	}
}

// TestOpenWritersRejectsDuplicate 覆盖重复项：
// MultiWriter 会把同一份日志写两遍，必须在构造期拒绝而不是静默去重。
func TestOpenWritersRejectsDuplicate(t *testing.T) {
	for _, paths := range [][]string{
		{"stdout", "stdout"},
		{"STDOUT", "stdout"}, // 忽略大小写后是同一个
		{"stderr", "STDERR"},
	} {
		if _, _, err := openWriters(config.LoggerConfig{OutputPaths: paths}); err == nil {
			t.Fatalf("%v 应被拒绝", paths)
		}
	}

	dir := t.TempDir()
	p := filepath.Join(dir, "a.log")
	if _, _, err := openWriters(config.LoggerConfig{OutputPaths: []string{p, p}}); err == nil {
		t.Fatal("同一个文件路径写两遍应被拒绝")
	}
	// 同一文件的另一种写法：./a.log 与 a.log 词法归一再比较
	alias := dir + string(os.PathSeparator) + "." + string(os.PathSeparator) + "a.log"
	if _, _, err := openWriters(config.LoggerConfig{OutputPaths: []string{p, alias}}); err == nil {
		t.Fatalf("%s 与 %s 是同一个文件，应被拒绝", p, alias)
	}
}

// TestOpenWritersClosesOnFailure 覆盖中途失败时的回滚：
// 第二个路径的父目录是普通文件，openWriters 必须返回错误而不是留下半开的状态。
func TestOpenWritersClosesOnFailure(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.log")
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := openWriters(config.LoggerConfig{
		OutputPaths: []string{p1, filepath.Join(blocker, "b.log")},
	})
	if err == nil {
		t.Fatal("第二个路径无法创建目录时应返回错误")
	}
}

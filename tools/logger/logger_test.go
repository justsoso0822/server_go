package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"app/config"
)

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

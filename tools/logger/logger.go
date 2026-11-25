// 提供基于项目配置创建日志记录器的辅助方法。
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"app/config"
)

// New 按配置构造 *slog.Logger，并返回关闭函数（关闭日志文件）。
//
// 输出由 slog 的 JSON/Text Handler 负责；时间、级别、caller 通过 ReplaceAttr
// 只把时间换成 time.DateTime 格式（字段名仍是 time）；level 和 source 都用 slog 默认。
func New(cfg config.LoggerConfig) (*slog.Logger, func() error, error) {
	var level slog.Level
	if strings.TrimSpace(cfg.Level) != "" {
		// slog.Level 实现了 encoding.TextUnmarshaler；解析失败直接返回错误，
		// 避免线上以错误级别静默运行。
		if err := level.UnmarshalText([]byte(strings.TrimSpace(cfg.Level))); err != nil {
			return nil, nil, fmt.Errorf("parse logger level %q: %w", cfg.Level, err)
		}
	}

	writers, closers, err := openWriters(cfg)
	if err != nil {
		return nil, nil, err
	}

	levelVar := new(slog.LevelVar)
	levelVar.Set(level)

	opts := &slog.HandlerOptions{
		Level:       levelVar,
		AddSource:   true,
		ReplaceAttr: replaceAttr,
	}

	var handler slog.Handler
	if loggerFormat(cfg) == "json" {
		handler = slog.NewJSONHandler(io.MultiWriter(writers...), opts)
	} else {
		handler = slog.NewTextHandler(io.MultiWriter(writers...), opts)
	}

	closeFn := func() error {
		var firstErr error
		for _, c := range closers {
			if err := c.Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	}

	return slog.New(handler), closeFn, nil
}

func loggerFormat(cfg config.LoggerConfig) string {
	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	switch format {
	case "json", "console":
		return format
	default:
		return "console"
	}
}

// replaceAttr 把 slog 的默认字段改成项目沿用的名字与格式：
//   - time   → time.DateTime 格式（字段名不变）
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	if a.Key == slog.TimeKey {
		return slog.String(slog.TimeKey, a.Value.Time().Format(time.DateTime))
	}
	return a
}

func openWriters(cfg config.LoggerConfig) ([]io.Writer, []*os.File, error) {
	var writers []io.Writer
	var closers []*os.File
	// seen 用来拒绝重复项：同一个 writer 出现两次，io.MultiWriter 会把
	// 每条日志写两遍（同一个文件则是两行交错），这属于配置错误。
	seen := make(map[string]struct{}, len(cfg.OutputPaths))

	for _, path := range cfg.OutputPaths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		// key 是「同一个输出」的判等依据：
		//   - 特殊值忽略大小写（配成 "STDOUT" 也该写标准输出，而不是建一个叫 STDOUT 的文件）；
		//   - 文件路径按词法归一（"./log/app.log" 与 "log/app.log" 是同一个文件），
		//     但保持原始大小写——Linux 区分大小写，整串 ToLower 会把两个真实文件折叠成一个。
		// 两类 key 加前缀区分，避免 "./stdout" 这种写法被误当成特殊值。
		key := "file:" + filepath.Clean(path)
		var special io.Writer
		switch strings.ToLower(path) {
		case "stdout":
			key, special = "stdout", os.Stdout
		case "stderr":
			key, special = "stderr", os.Stderr
		}

		if _, dup := seen[key]; dup {
			closeFiles(closers)
			return nil, nil, fmt.Errorf("duplicate log output path %q in outputPaths", path)
		}
		seen[key] = struct{}{}

		if special != nil {
			writers = append(writers, special)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			closeFiles(closers)
			return nil, nil, fmt.Errorf("create log dir for %s: %w", path, err)
		}
		// O_TRUNC 是按文件生效的：每个路径各自截断，不存在"清多次"的问题，
		// 所以这里只在打开时判断，绝不能把 ClearOnStart 改成 false。
		flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
		if cfg.ClearOnStart {
			flags |= os.O_TRUNC
		}
		file, err := os.OpenFile(path, flags, 0o644)
		if err != nil {
			// 已经打开的 N-1 个文件要关掉，否则构造失败会泄漏 fd。
			closeFiles(closers)
			return nil, nil, fmt.Errorf("open log file %s: %w", path, err)
		}
		writers = append(writers, file)
		closers = append(closers, file)
	}

	if len(writers) == 0 {
		// 什么都没配时至少保证有输出，避免日志静默丢失。
		writers = append(writers, os.Stdout)
	}
	return writers, closers, nil
}

// closeFiles 关闭构造日志过程中已打开的文件。
// 只用于失败回滚：成功时这些句柄由 New 返回的关闭函数统一关。
func closeFiles(files []*os.File) {
	for _, f := range files {
		_ = f.Close()
	}
}

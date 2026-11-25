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
// 保持与旧 zap 日志一致的字段名（ts / level / caller）。
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
//   - time   → ts（time.DateTime 格式）
//   - level  → 小写（error/warn/info/debug）
//   - source → caller（目录/文件:行号）
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	switch a.Key {
	case slog.TimeKey:
		return slog.String("ts", a.Value.Time().Format(time.DateTime))
	case slog.LevelKey:
		if lv, ok := a.Value.Any().(slog.Level); ok {
			return slog.String("level", strings.ToLower(lv.String()))
		}
		return a
	case slog.SourceKey:
		if src, ok := a.Value.Any().(*slog.Source); ok {
			return slog.String("caller", shortSource(src))
		}
		return a
	default:
		return a
	}
}

// shortSource 把绝对路径裁成「最后两级 + 行号」，例如 middleware/access_log.go:42。
func shortSource(src *slog.Source) string {
	if src == nil {
		return ""
	}
	file := filepath.ToSlash(src.File)
	if parts := strings.Split(file, "/"); len(parts) >= 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}
	return fmt.Sprintf("%s:%d", file, src.Line)
}

func openWriters(cfg config.LoggerConfig) ([]io.Writer, []*os.File, error) {
	var writers []io.Writer
	var closers []*os.File

	for _, path := range cfg.OutputPaths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		switch path {
		case "stdout":
			writers = append(writers, os.Stdout)
			continue
		case "stderr":
			writers = append(writers, os.Stderr)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, nil, fmt.Errorf("create log dir for %s: %w", path, err)
		}
		flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
		if cfg.ClearOnStart {
			flags |= os.O_TRUNC
			cfg.ClearOnStart = false // 只清一次
		}
		file, err := os.OpenFile(path, flags, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("open log file %s: %w", path, err)
		}
		writers = append(writers, file)
		closers = append(closers, file)
	}

	if cfg.Stdout {
		writers = append(writers, os.Stdout)
	}
	if len(writers) == 0 {
		// 什么都没配时至少保证有输出，避免日志静默丢失。
		writers = append(writers, os.Stdout)
	}
	return writers, closers, nil
}

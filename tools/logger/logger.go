package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server_go/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LoggerConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("parse logger level %q: %w", cfg.Level, err)
	}

	var zapCfg zap.Config
	switch loggerFormat(cfg) {
	case "json":
		zapCfg = zap.NewProductionConfig()
	default:
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeTime = ginStyleTimeEncoder
		if cfg.StdoutColorDisabled {
			zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.OutputPaths = append([]string(nil), cfg.OutputPaths...)
	if cfg.Stdout {
		zapCfg.OutputPaths = append(zapCfg.OutputPaths, "stdout")
	}
	if err := ensureLogDirs(zapCfg.OutputPaths); err != nil {
		return nil, err
	}
	if cfg.ClearOnStart {
		if err := clearLogFiles(zapCfg.OutputPaths); err != nil {
			return nil, err
		}
	}

	log, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}
	return log, nil
}

func loggerFormat(cfg config.LoggerConfig) string {
	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	if format != "" {
		return format
	}
	if cfg.StdoutColorDisabled {
		return "json"
	}
	return "console"
}

func ginStyleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 - 15:04:05"))
}

func ensureLogDirs(paths []string) error {
	for _, path := range paths {
		if path == "" || path == "stdout" || path == "stderr" {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("create log dir for %s: %w", path, err)
		}
	}
	return nil
}

func clearLogFiles(paths []string) error {
	for _, path := range paths {
		if path == "" || path == "stdout" || path == "stderr" {
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("clear log file %s: %w", path, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("close cleared log file %s: %w", path, err)
		}
	}
	return nil
}

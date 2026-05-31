package logger

import (
	"fmt"
	"strings"

	"server_gin/config"

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
		if cfg.StdoutColorDisabled {
			zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)
	if cfg.Stdout {
		zapCfg.OutputPaths = []string{"stdout"}
	} else {
		zapCfg.OutputPaths = nil
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

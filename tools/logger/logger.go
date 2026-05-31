package logger

import (
	"strings"

	"server_gin/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LoggerConfig) *zap.Logger {
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level))

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
		log = zap.NewNop()
	}
	return log
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

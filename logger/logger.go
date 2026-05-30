package logger

import (
	"server_gin/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LoggerConfig) *zap.Logger {
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level))

	var zapCfg zap.Config
	if cfg.StdoutColorDisabled {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.OutputPaths = []string{"stdout"}

	log, err := zapCfg.Build()
	if err != nil {
		log = zap.NewNop()
	}
	return log
}

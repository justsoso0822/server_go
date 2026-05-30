package bootstrap

import (
	"fmt"

	"server_gin/config"
	"server_gin/tools/autodb"
	"server_gin/tools/logger"

	"go.uber.org/zap"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.Logger)

	if err := autodb.Init(cfg); err != nil {
		_ = log.Sync()
		return nil, fmt.Errorf("init db: %w", err)
	}

	return &App{
		Config: cfg,
		Logger: log,
	}, nil
}

func (app *App) Close() error {
	var err error
	if e := autodb.Close(); e != nil {
		err = e
	}
	if app != nil && app.Logger != nil {
		if e := app.Logger.Sync(); e != nil && err == nil {
			err = e
		}
	}
	return err
}

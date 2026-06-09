// 启动包负责应用的初始化和资源管理，包括配置加载、日志初始化和数据库连接等。
package bootstrap

import (
	"fmt"

	"server_go/config"
	"server_go/tools/autodb"
	"server_go/tools/logger"

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

	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	zap.ReplaceGlobals(log)

	if err := autodb.Init(cfg, log); err != nil {
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

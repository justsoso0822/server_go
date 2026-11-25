// 启动包负责应用的初始化和资源管理，包括配置加载、日志初始化和数据库连接等。
package bootstrap

import (
	"fmt"
	"log/slog"

	"app/config"
	"app/tools/autodb"
	"app/tools/logger"
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	closeLogger func() error
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log, closeLog, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	// 设为全局 logger：service/middleware 里 slog.Default() 拿到的就是它。
	slog.SetDefault(log)

	if err := autodb.Init(cfg, log); err != nil {
		_ = closeLog()
		return nil, fmt.Errorf("init db: %w", err)
	}

	return &App{
		Config:      cfg,
		Logger:      log,
		closeLogger: closeLog,
	}, nil
}

func (app *App) Close() error {
	var err error
	if e := autodb.Close(); e != nil {
		err = e
	}
	if app != nil && app.closeLogger != nil {
		if e := app.closeLogger(); e != nil && err == nil {
			err = e
		}
	}
	return err
}

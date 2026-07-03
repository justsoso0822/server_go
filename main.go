// Package main 是程序入口，负责启动 HTTP 服务器并处理优雅关闭。
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"server_go/bootstrap"
	"server_go/router"

	"go.uber.org/zap"
)

// main 是程序主入口，启动失败时将错误写入 stderr 并以非零状态码退出。
func main() {
	if err := run(); err != nil {
		writeStartupError(err)
		os.Exit(1)
	}
}

// writeStartupError 将启动错误输出到 stderr。
// 本地环境（local）输出纯文本，其他环境输出 JSON 格式，便于日志系统采集。
func writeStartupError(err error) {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env == "" || env == "local" {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	_ = json.NewEncoder(os.Stderr).Encode(map[string]any{
		"level": "error",
		"msg":   "startup failed",
		"error": err.Error(),
	})
}

// run 负责初始化应用、启动 HTTP 服务器，并监听系统信号以实现优雅关闭。
func run() error {
	// 初始化应用（加载配置、数据库连接、日志等）
	app, err := bootstrap.New()
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()

	// 配置 HTTP 服务器，设置合理的超时时间防止慢请求耗尽资源
	srv := &http.Server{
		Addr:              app.Config.Server.Address,
		Handler:           router.Setup(app),
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 在 goroutine 中启动服务器，避免阻塞主流程
	serverErr := make(chan error, 1)
	go func() {
		app.Logger.Info("http server started", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	// 监听 SIGINT（Ctrl+C）和 SIGTERM（kill）信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	// 等待服务器异常或收到退出信号
	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		app.Logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	// 优雅关闭：等待正在处理的请求完成，超时时间由配置决定
	timeout := time.Duration(app.Config.Server.GracefulTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return <-serverErr
}

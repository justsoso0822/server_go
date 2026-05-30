package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server_gin/bootstrap"
	"server_gin/router"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	app, err := bootstrap.New()
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()

	srv := &http.Server{
		Addr:    app.Config.Server.Address,
		Handler: router.Setup(app),
	}

	serverErr := make(chan error, 1)
	go func() {
		app.Logger.Info("http server started", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		app.Logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	timeout := time.Duration(app.Config.Server.GracefulTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return <-serverErr
}

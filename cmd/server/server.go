package server

import (
	"context"
	"gin-boilerplate/config"
	"gin-boilerplate/database/sql/postgres"
	"gin-boilerplate/pkg/logger"
	"gin-boilerplate/pkg/storage/local"
	"gin-boilerplate/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap/zapcore"
)

func StartServer() {
	// Load config
	config.InitConfig()
	cfg := config.GetConfig()

	logger.InitLogger(cfg)
	defer logger.Sync()

	ctx := context.Background()

	// mongodb.InitClient(ctx)
	postgres.InitClient(ctx, cfg)
	local.InitLocalStorage()
	validator.InitValidator()

	// Register all dependencies
	Register()

	r := InitHandler()

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "Server startup failed:", zapcore.Field{Key: "error", Type: zapcore.StringType, String: err.Error()})
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info(ctx, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(ctx, "Server Shutdown:", zapcore.Field{Key: "error", Type: zapcore.StringType, String: err.Error()})
	}
	<-ctx.Done()
	logger.Info(ctx, "Server shutdown gracefully.")
}

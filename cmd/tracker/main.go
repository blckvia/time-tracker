package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"time-tracker/internal/app"
)

func main() {
	ctx := context.Background()

	logger := zap.Must(zap.NewProduction())
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Error syncing logger", zap.Error(err))
		}
	}(logger)

	tracker := app.NewApp(ctx, logger)

	go func() {
		if err := tracker.Run(); err != nil {
			logger.Fatal("Error running server", zap.Error(err))
		}
	}()

	tracker.Logger.Info("Server started", zap.String("address", tracker.Server.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	tracker.Logger.Info("Shutting down server")

	if err := tracker.Shutdown(ctx); err != nil {
		logger.Fatal("Error shutting down server", zap.Error(err))
	}
}

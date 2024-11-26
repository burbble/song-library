package main

import (
	"context"
	"fmt"
	"song-library/internal/app"
	"song-library/internal/config"
	"song-library/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	log := logger.New(cfg.Log.Level)
	ctx := context.Background()

	log.Info(ctx, "Starting application")

	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatal("Application creation error", zap.Error(err))
	}
	log.Info(ctx, "Application successfully created")

	if err := application.Run(); err != nil {
		log.Fatal("Application runtime error", zap.Error(err))
	}
}

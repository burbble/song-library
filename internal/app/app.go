package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "song-library/docs"
	"song-library/internal/application/usecase"
	"song-library/internal/config"
	"song-library/internal/infrastructure/database"
	"song-library/internal/infrastructure/persistence/postgres"
	"song-library/internal/interfaces/http/handler"
	"song-library/pkg/logger"
)

type App struct {
	config *config.Config
	router *gin.Engine
	logger *logger.Logger
	db     *database.Database
}

func New(cfg *config.Config, logger *logger.Logger) (*App, error) {
	ctx := context.Background()
	
	logger.Info(ctx, "Starting application initialization")
	
	db, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Error(ctx, "Database initialization error", zap.Error(err))
		return nil, fmt.Errorf("database initialization error: %w", err)
	}
	logger.Debug(ctx, "Database successfully initialized")

	app := &App{
		config: cfg,
		router: gin.Default(),
		logger: logger,
		db:     db,
	}

	app.setupRoutes(logger)
	logger.Info(ctx, "Routes successfully configured")
	
	return app, nil
}

func (a *App) setupRoutes(logger *logger.Logger) {
	songRepo := postgres.NewSongRepository(a.db.GetDB(), logger)
	songUseCase := usecase.NewSongUseCase(songRepo, a.config)
	songHandler := handler.NewSongHandler(*songUseCase, logger)

	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := a.router.Group("/api/v1")
	{
		songs := v1.Group("/songs")
		{
			songs.POST("", songHandler.Create)
			songs.GET("", songHandler.List)
			songs.GET("/:id", songHandler.Get)
			songs.PUT("/:id", songHandler.Update)
			songs.DELETE("/:id", songHandler.Delete)
			songs.GET("/:id/text", songHandler.GetSongText)
		}
	}
}

func (a *App) Run() error {
	ctx := context.Background()
	
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port),
		Handler: a.router,
	}

	errChan := make(chan error, 1)

	go func() {
		a.logger.Info(ctx, "Starting server", 
			zap.String("host", a.config.Server.Host),
			zap.String("port", a.config.Server.Port))
			
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error(ctx, "Server start error", zap.Error(err))
			errChan <- fmt.Errorf("server start error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case sig := <-quit:
		a.logger.Info(ctx, "Received termination signal", zap.String("signal", sig.String()))
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			a.logger.Error(ctx, "Graceful shutdown error", zap.Error(err))
			return fmt.Errorf("graceful shutdown error: %w", err)
		}

		a.logger.Info(ctx, "Server successfully stopped")
	}

	return nil
}

func (a *App) Close() error {
	return a.db.Close()
}

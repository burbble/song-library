package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"song-library/internal/config"
	"song-library/pkg/logger"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(cfg *config.Config, logger *logger.Logger) (*Database, error) {
	ctx := context.Background()
	
	logger.Debug(ctx, "Starting database connection", 
		zap.String("host", cfg.Database.Host),
		zap.String("port", cfg.Database.Port),
		zap.String("dbname", cfg.Database.DBName))

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error(ctx, "Database connection error", zap.Error(err))
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	if err = db.Ping(); err != nil {
		logger.Error(ctx, "Database ping error", zap.Error(err))
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	logger.Info(ctx, "Successfully connected to database")
	return &Database{db: db}, nil
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) Close() error {
	return d.db.Close()
}

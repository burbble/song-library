package config

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"song-library/pkg/logger"
	"go.uber.org/zap"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
	Log      struct {
		Level string `env:"LOG_LEVEL" envDefault:"info"`
	}
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type APIConfig struct {
	MusicInfoURL string
}

func LoadConfig() (*Config, error) {
	log := logger.New("debug")
	ctx := context.Background()
	
	log.Debug(ctx, "Начало загрузки конфигурации")

	if err := godotenv.Load(); err != nil {
		log.Warn(ctx, "Файл .env не найден, используются значения по умолчанию", zap.Error(err))
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "song_library_user"),
			Password: getEnv("DB_PASSWORD", "song_library_password"),
			DBName:   getEnv("DB_NAME", "song_library_db"),
		},
		API: APIConfig{
			MusicInfoURL: getEnv("MUSIC_INFO_API_URL", "http://localhost:8081"),
		},
	}

	log.Info(ctx, "Конфигурация успешно загружена", 
		zap.String("server_host", cfg.Server.Host),
		zap.String("server_port", cfg.Server.Port),
		zap.String("db_host", cfg.Database.Host),
		zap.String("db_name", cfg.Database.DBName))

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresUser string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresPort string `env:"POSTGRES_PORT"`
	PostgresDBName string `env:"POSTGRES_DB"`
	DatabaseURL string `env:"DATABASE_URL"`
	JWTSecret   string `env:"JWT_SECRET"`
	JWTExpiration int    `env:"JWT_EXPIRATION"`
	Env string `env:"ENV"`
}

func LoadConfig() (*Config, error) {
	slog.Debug("load .env")
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	cfg.PostgresUser = os.Getenv("POSTGRES_USER")
	cfg.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	cfg.PostgresPort = os.Getenv("POSTGRES_PORT")
	cfg.PostgresDBName = os.Getenv("POSTGRES_DB")
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.JWTExpiration = 72 // default to 72 hours if not set
	if expiration := os.Getenv("JWT_EXPIRATION"); expiration != "" {
		value, err := strconv.Atoi(expiration)
		if err != nil {
			slog.Warn("invalid jwt expiration duration in .env, using default 72hours")
		} else {
			cfg.JWTExpiration = value
		}
	}
	slog.Info("successfully loaded .env")
	return &cfg, nil
}
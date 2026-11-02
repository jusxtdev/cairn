package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresUser string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresPort string `env:"POSTGRES_PORT,required"`
	PostgresDBName string `env:"POSTGRES_DB,required"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	JWTSecret   string `env:"JWT_SECRET"`
	JWTExpiration int    `env:"JWT_EXPIRATION"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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
		if value, err := strconv.Atoi(expiration); err == nil {
			cfg.JWTExpiration = value
		}
	}
	return &cfg, nil
}
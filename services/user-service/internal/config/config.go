package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	RedisAddr   string
	ServerPort  string
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ServerPort:  getEnv("SERVER_PORT", "8081"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

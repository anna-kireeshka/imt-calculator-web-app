package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DataBaseUrl   string
	ListenAddress string
	TgBotToken    string
}

func LoadConfig() (*Config, error) {
	initEnv()

	var cfg = &Config{
		DataBaseUrl:   getEnv("DATABASE_URL", ""),
		ListenAddress: getEnv("LISTEN_ADDRESS", ":8080"),
		TgBotToken:    getEnv("TELEGRAM_API_TOKEN", ""),
	}

	if cfg.DataBaseUrl == "" {
		return nil, errors.New("DATABASE_URL must be set")
	}

	if cfg.TgBotToken == "" {
		return nil, errors.New("TELEGRAM_API_TOKEN must be set")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

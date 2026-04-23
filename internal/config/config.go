package config

import (
	"log"
	"os"
)

var Cfg *Config

type Config struct {
	DbURL   string
	AppPort string
	Env     string
}

func LoadConfig() {
	dbURL := os.Getenv("DB_DSN")
	if dbURL == "" {
		log.Fatal("DB_DSN не установлена в .env")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	Cfg = &Config{
		DbURL:   dbURL,
		AppPort: appPort,
		Env:     os.Getenv("ENV"),
	}
}

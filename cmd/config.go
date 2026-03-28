package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	ConfigBaseDomain = "BASE_DOMAIN"
	ConfigAddress    = "ADDRESS"

	ConfigPostgresHost     = "POSTGRES_HOST"
	ConfigPostgresDatabase = "POSTGRES_DATABASE"
	ConfigPostgresUser     = "POSTGRES_USER"
	ConfigPostgresPassword = "POSTGRES_PASSWORD"
)

var configs = []string{
	ConfigBaseDomain,
	ConfigAddress,

	ConfigPostgresHost,
	ConfigPostgresDatabase,
	ConfigPostgresUser,
	ConfigPostgresPassword,
}

func checkConfigs() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("Error loading .env file: %v", err)
	}

	for _, config := range configs {
		if os.Getenv(config) == "" {
			log.Fatalf("Config %s is not set", config)
		}
	}
}

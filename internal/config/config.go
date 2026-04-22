package config

import (
	"os"
)

type Config struct {
	Port     string
	DBDriver string
	DBURL    string
	DBPath   string
	LogLevel string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		DBDriver: getEnv("DB_DRIVER", "sqlite3"),
		DBURL:    getEnv("DB_URL", ""),
		DBPath:   getEnv("DB_PATH", "./books.db"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

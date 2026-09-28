package config

import "os"

type Config struct {
	AppPort  string
	LogLevel string

	PostgresDSN string

	RedisURL string
}

func Load() Config {
	return Config{
		AppPort:     getEnv("APP_PORT", "3000"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		PostgresDSN: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/playground?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

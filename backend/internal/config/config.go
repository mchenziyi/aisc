package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTExpiration      time.Duration
	ServerPort         string
	CORSAllowedOrigins string
	LogLevel           string
	DBMaxConns         int
	DBMinConns         int
}

func Load() *Config {
	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production-at-least-32-chars!!")
	if jwtSecret == "change-me-in-production-at-least-32-chars!!" {
		log.Println("WARNING: Using default JWT secret. This is insecure for production. Set JWT_SECRET environment variable.")
	}

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable"),
		JWTSecret:          jwtSecret,
		JWTExpiration:      getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DBMaxConns:         getIntEnv("DB_MAX_CONNS", 25),
		DBMinConns:         getIntEnv("DB_MIN_CONNS", 5),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
	}
	return fallback
}

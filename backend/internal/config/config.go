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
	RunMigrations      bool
	DBMaxConns         int
	DBMinConns         int
}

func Load() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is required")
	}

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable"),
		JWTSecret:          jwtSecret,
		JWTExpiration:      getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		RunMigrations:      getEnv("RUN_MIGRATIONS", "false") == "true",
		DBMaxConns:         getIntEnv("DB_MAX_CONNS", 25),
		DBMinConns:         getIntEnv("DB_MIN_CONNS", 5),
	}
}

// LoadForMigrate loads only database-related configuration for migration commands.
// It does NOT require JWT_SECRET, allowing migrations to run without it.
func LoadForMigrate() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable"),
		RunMigrations: true,
		DBMaxConns:    getIntEnv("DB_MAX_CONNS", 25),
		DBMinConns:    getIntEnv("DB_MIN_CONNS", 5),
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
		log.Printf("warning: invalid duration value for %s (%q), using default %v", key, v, fallback)
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
		log.Printf("warning: invalid integer value for %s (%q), using default %d", key, v, fallback)
	}
	return fallback
}

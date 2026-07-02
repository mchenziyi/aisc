package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application.
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

// Load reads configuration from environment variables.
func Load() *Config {
	cfg := &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/todoapp"),
		JWTSecret:          getEnvRequired("JWT_SECRET"),
		JWTExpiration:      getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DBMaxConns:         getIntEnv("DB_MAX_CONNS", 25),
		DBMinConns:         getIntEnv("DB_MIN_CONNS", 5),
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvRequired returns the environment variable value or panics if empty.
func getEnvRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("required environment variable " + key + " is not set")
	}
	return val
}

func getIntEnv(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

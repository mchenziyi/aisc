package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Log levels
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var currentLogLevel = LevelInfo

func init() {
	lvl := os.Getenv("LOG_LEVEL")
	if lvl != "" {
		lvl = strings.ToLower(lvl)
		switch lvl {
		case LevelDebug, LevelInfo, LevelWarn, LevelError:
			currentLogLevel = lvl
		default:
			log.Printf("warning: invalid LOG_LEVEL %q, using default 'info'", lvl)
		}
	}
}

// shouldLog returns true if the given level should be logged based on current log level.
func shouldLog(level string) bool {
	levels := []string{LevelDebug, LevelInfo, LevelWarn, LevelError}
	currentIdx := 0
	levelIdx := 0
	for i, l := range levels {
		if l == currentLogLevel {
			currentIdx = i
		}
		if l == level {
			levelIdx = i
		}
	}
	return levelIdx >= currentIdx
}

// LoggerMiddleware creates a request logging middleware.
// It generates a unique request_id for each request and logs method, path, status, and duration.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate and set request_id
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log after request is processed
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Determine log level based on status code
		logLevel := LevelInfo
		if status >= 500 {
			logLevel = LevelError
		} else if status >= 400 {
			logLevel = LevelWarn
		}

		if shouldLog(logLevel) {
			log.Printf("[%s] %s %s %d %v", requestID, method, path, status, duration)
		}
	}
}

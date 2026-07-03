package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Log levels
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// shouldLog returns true if the given level should be logged based on current log level.
func shouldLog(level string, currentLogLevel string) bool {
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
func LoggerMiddleware(logLevel string) gin.HandlerFunc {
	// Normalize log level
	logLevel = strings.ToLower(logLevel)
	switch logLevel {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		// valid
	case "":
		logLevel = LevelInfo
	default:
		log.Printf("warning: invalid LOG_LEVEL %q, using default 'info'", logLevel)
		logLevel = LevelInfo
	}

	return func(c *gin.Context) {
		start := time.Now()

		// Generate and set request_id
		requestID := generateRequestID()
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
		msgLevel := LevelInfo
		if status >= 500 {
			msgLevel = LevelError
		} else if status >= 400 {
			msgLevel = LevelWarn
		}

		if shouldLog(msgLevel, logLevel) {
			log.Printf("[%s] %s %s %d %v", requestID, method, path, status, duration)
		}
	}
}

// generateRequestID generates a unique request ID using crypto/rand.
func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format("20060102150405.000000")
	}
	return hex.EncodeToString(b)
}

// GetRequestID retrieves the request_id from the Gin context.
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}

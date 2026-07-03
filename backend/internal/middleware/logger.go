package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
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

// logLevelNumeric maps log level strings to numeric values for efficient comparison.
var logLevelNumeric = map[string]int{
	LevelDebug: 0,
	LevelInfo:  1,
	LevelWarn:  2,
	LevelError: 3,
}

// messageLevel computes the log level for a given HTTP status code.
func messageLevel(status int) string {
	switch {
	case status >= 500:
		return LevelError
	case status >= 400:
		return LevelWarn
	default:
		return LevelInfo
	}
}

// LoggerMiddleware creates a request logging middleware.
// It generates a unique request_id for each request and logs method, path, status, and duration.
func LoggerMiddleware(logLevel string) gin.HandlerFunc {
	// Normalize log level and get its numeric value for efficient comparison
	effectiveLevel := LevelInfo
	if n, ok := logLevelNumeric[logLevel]; ok {
		effectiveLevel = logLevel
		_ = n
	} else if logLevel != "" {
		log.Printf("warning: invalid LOG_LEVEL %q, using default 'info'", logLevel)
	}

	currentLevelNum := logLevelNumeric[effectiveLevel]

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
		msgLevel := messageLevel(status)
		if msgLevelNum, ok := logLevelNumeric[msgLevel]; ok && msgLevelNum >= currentLevelNum {
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

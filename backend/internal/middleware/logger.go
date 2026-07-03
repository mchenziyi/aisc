package middleware

import (
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
func LoggerMiddleware(logLevel string) gin.HandlerFunc {
	// Normalize log level
	effectiveLevel := LevelInfo
	if _, ok := logLevelNumeric[logLevel]; ok {
		effectiveLevel = logLevel
	} else if logLevel != "" {
		log.Printf("warning: invalid LOG_LEVEL %q, using default 'info'", logLevel)
	}

	currentLevelNum := logLevelNumeric[effectiveLevel]

	return func(c *gin.Context) {
		start := time.Now()

		// Get request_id from context (set by RequestIDMiddleware)
		requestID := c.GetString("request_id")

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

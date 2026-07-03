package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a CORS middleware with configurable allowed origins.
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")
	// Build a set for O(1) lookup
	originSet := make(map[string]bool, len(origins))
	hasWildcard := false
	for i := range origins {
		o := strings.TrimSpace(origins[i])
		if o == "*" {
			hasWildcard = true
		}
		if o != "" {
			originSet[o] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed using map lookup (O(1))
		allowed := hasWildcard || originSet[origin]

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			c.Header("Access-Control-Expose-Headers", "X-Request-ID")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

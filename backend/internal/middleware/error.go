package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

// ErrorMiddleware is a Gin middleware that catches errors added via c.Error()
// and formats them into the standard error response format.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only process if there are errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error (most specific)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// Determine the request ID
		requestID := getRequestID(c)

		// If it's an AppError, use its fields
		if appErr, ok := err.Err.(*apperrors.AppError); ok {
			appErr.RequestID = requestID
			c.JSON(appErr.Code, appErr)
			return
		}

		// Generic fallback
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":       500,
			"error_code": apperrors.ErrorCodeInternal,
			"message":    "internal server error",
			"request_id": requestID,
		})
	}
}

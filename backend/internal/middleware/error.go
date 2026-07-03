package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
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

		// If it's an AppError, use its fields
		if appErr, ok := err.Err.(*apperrors.AppError); ok {
			resp := model.NewErrorResponse(appErr.Code, appErr.Message, appErr.FieldErrors)
			c.JSON(appErr.HTTPCode, resp)
			return
		}

		// Generic fallback
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			apperrors.CodeInternal,
			"服务器内部错误，请稍后重试",
			nil,
		))
	}
}

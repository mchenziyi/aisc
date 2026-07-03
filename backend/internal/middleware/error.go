package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

// ErrorMiddleware catches errors added via c.Error() and formats them
// using the standard error response format.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		if err == nil {
			return
		}

		// If it's an AppError, use RespondError
		if appErr, ok := err.Err.(*apperrors.AppError); ok {
			apperrors.RespondError(c, appErr)
			c.Abort()
			return
		}

		// Generic fallback
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "服务器内部错误，请稍后重试",
		})
	}
}

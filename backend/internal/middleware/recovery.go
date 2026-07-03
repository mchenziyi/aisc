package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

// RecoveryMiddleware returns a middleware that recovers from panics
// and responds with a unified ErrorResponse format instead of the
// default plain-text panic output.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{
					Code:    apperrors.CodeInternal,
					Message: "服务器内部错误，请稍后重试",
				})
			}
		}()
		c.Next()
	}
}

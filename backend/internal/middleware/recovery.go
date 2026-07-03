package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"todo-api/internal/errors"
	"todo-api/internal/model"
)

// RecoveryMiddleware returns a middleware that recovers from panics
// and responds with a unified error format.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				requestID := c.GetString("request_id")
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse{
					ErrorCode: errors.CodeInternal,
					Message:   "服务器内部错误，请稍后重试",
					RequestID: requestID,
				})
			}
		}()
		c.Next()
	}
}

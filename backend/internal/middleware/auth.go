package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"todo-api/internal/auth"
	apperrors "todo-api/internal/errors"
)

// AuthMiddleware creates a JWT authentication middleware.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apperrors.RespondError(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			apperrors.RespondError(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims := &auth.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				apperrors.RespondError(c, apperrors.ErrTokenExpired)
			} else {
				apperrors.RespondError(c, apperrors.ErrInvalidToken)
			}
			c.Abort()
			return
		}

		if !token.Valid {
			apperrors.RespondError(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

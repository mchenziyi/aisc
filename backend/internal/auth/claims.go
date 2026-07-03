package auth

import "github.com/golang-jwt/jwt/v5"

// UserClaims represents the custom JWT claims including user_id as int64.
type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

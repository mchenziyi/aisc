package auth

import "github.com/golang-jwt/jwt/v5"

// UserClaims represents the custom JWT claims containing user_id.
type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

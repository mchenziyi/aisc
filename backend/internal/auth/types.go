package auth

import (
	"time"
)

// ─── Request DTOs ─────────────────────────────────────────────

// RegisterRequest represents the registration request body.
// Username: 3-32 chars, only letters and digits.
// Password: 8-128 chars.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the refresh token request body.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ─── Internal DTOs ────────────────────────────────────────────

// User represents a user record from the database.
type User struct {
	ID                  int64      `json:"id"`
	Username            string     `json:"username"`
	Password            string     `json:"-"`
	RefreshTokenHash    *string    `json:"-"`
	RefreshTokenExpires *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

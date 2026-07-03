package auth

import (
	"time"
)

// ─── Request DTOs ─────────────────────────────────────────────

// RegisterRequest represents the registration request body.
// Username: 3-20 chars, letters, digits, and underscores allowed.
// Password: 8-128 chars, must contain both letters and digits.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ─── Database Model ───────────────────────────────────────────

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

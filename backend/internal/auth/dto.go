package auth

import "time"

// RegisterRequest represents the register request body.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response.
type LoginResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

// UserPublic represents the public user information.
type UserPublic struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// User represents a full user record from the database.
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

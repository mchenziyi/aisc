package auth

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse represents the registration success response.
type RegisterResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login success response.
type LoginResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

// UserPublic represents the public user information.
type UserPublic struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

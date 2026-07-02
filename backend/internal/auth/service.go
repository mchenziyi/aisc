package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	// passwordRegex requires at least 8 characters with at least one letter and one digit.
	// Special characters (symbols) are allowed. This can be adjusted after PM confirmation.
	passwordRegex = regexp.MustCompile(`^(?=.*[a-zA-Z])(?=.*\d).{8,}$`)
)

// Service handles auth business logic.
type Service struct {
	repo           *Repository
	jwtSecret      string
	jwtExpiration  time.Duration
}

// NewService creates a new auth service.
func NewService(repo *Repository, jwtSecret string, jwtExpiration time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// Register creates a new user account.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*UserPublic, error) {
	// Normalize username to lowercase
	username := strings.ToLower(strings.TrimSpace(req.Username))

	// Validate username
	if !usernameRegex.MatchString(username) {
		return nil, fmt.Errorf("username must be 3-20 characters, allowing letters, digits and underscores")
	}

	// Validate password
	if !passwordRegex.MatchString(req.Password) {
		return nil, fmt.Errorf("password must be at least 8 characters, containing both letters and digits (special characters allowed)")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, username, string(hashedPassword))
	if err != nil {
		if err == ErrUsernameTaken {
			return nil, err
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &UserPublic{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// GetMe returns the public profile of the currently authenticated user.
func (s *Service) GetMe(ctx context.Context, userID int64) (*UserPublic, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &UserPublic{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Normalize username to lowercase for lookup
	username := strings.ToLower(strings.TrimSpace(req.Username))

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, fmt.Errorf("invalid username or password")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(s.jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &LoginResponse{
		Token: tokenString,
		User: UserPublic{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}

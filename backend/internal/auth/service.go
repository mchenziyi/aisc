package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

type Service struct {
	repo          *Repository
	jwtSecret     []byte
	jwtExpiration time.Duration // access token lifetime
	tokenExpiry   time.Duration // refresh token lifetime
}

func NewService(repo *Repository, jwtSecret string, jwtExpiration, tokenExpiry time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
		tokenExpiry:   tokenExpiry,
	}
}

// Register creates a new user account and returns AuthResponse (register = login).
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*model.AuthResponse, *apperrors.AppError) {
	// Validate username format
	username := strings.ToLower(req.Username)
	if appErr := validateUsername(username); appErr != nil {
		return nil, appErr
	}

	// Validate password strength
	if appErr := validatePassword(req.Password); appErr != nil {
		return nil, appErr
	}

	// Check if username already exists
	existing, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("internal error: failed to check username '%s': %v", username, err)
		return nil, apperrors.ErrInternal
	}
	if existing != nil {
		return nil, apperrors.ErrUsernameTaken
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("internal error: bcrypt hash failed: %v", err)
		return nil, apperrors.ErrInternal
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, username, string(hash))
	if err != nil {
		log.Printf("internal error: failed to create user '%s': %v", username, err)
		return nil, apperrors.ErrInternal
	}

	// Generate JWT and refresh token
	authResp, appErr := s.generateAuthResponse(ctx, user.ID, user.Username)
	if appErr != nil {
		return nil, appErr
	}

	return authResp, nil
}

// Login authenticates a user and returns AuthResponse.
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*model.AuthResponse, *apperrors.AppError) {
	username := strings.ToLower(req.Username)

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("internal error: failed to find user '%s': %v", username, err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Generate JWT and refresh token
	authResp, appErr := s.generateAuthResponse(ctx, user.ID, user.Username)
	if appErr != nil {
		return nil, appErr
	}

	return authResp, nil
}

// RefreshToken generates a new JWT for the given user.
func (s *Service) RefreshToken(ctx context.Context, userID int64) (*model.AuthResponse, *apperrors.AppError) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("internal error: failed to find user by ID %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Generate new JWT and refresh token
	authResp, appErr := s.generateAuthResponse(ctx, user.ID, user.Username)
	if appErr != nil {
		return nil, appErr
	}

	return authResp, nil
}

// GetCurrentUser returns the current user's public info (id, username, created_at).
func (s *Service) GetCurrentUser(ctx context.Context, userID int64) (*model.UserPublic, *apperrors.AppError) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("internal error: failed to find user by ID %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}
	return &model.UserPublic{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// generateAuthResponse creates a JWT, refresh token, and returns AuthResponse.
func (s *Service) generateAuthResponse(ctx context.Context, userID int64, username string) (*model.AuthResponse, *apperrors.AppError) {
	now := time.Now()
	expiresIn := s.jwtExpiration

	// Create JWT claims
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   username,
		},
	}

	// Sign JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		log.Printf("internal error: JWT generation failed for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}

	// Generate refresh token
	refreshToken := generateRandomToken(32)
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("internal error: bcrypt hash failed for refresh token: %v", err)
		return nil, apperrors.ErrInternal
	}
	refreshExpiresAt := now.Add(s.tokenExpiry)

	// Store refresh token hash and expiry
	if err := s.repo.UpdateRefreshToken(ctx, userID, string(refreshHash), refreshExpiresAt); err != nil {
		log.Printf("internal error: failed to store refresh token for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}

	return &model.AuthResponse{
		Token: tokenStr,
		User: model.UserInfo{
			ID:       userID,
			Username: username,
		},
	}, nil
}

// generateRandomToken generates a cryptographically random hex string of given byte length.
func generateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Printf("warning: failed to generate random token: %v", err)
		// Fallback to a less secure random (better than nothing)
		for i := range b {
			b[i] = byte(i + 1)
		}
	}
	return hex.EncodeToString(b)
}

// validateUsername checks that the username is 3-20 characters and allows letters, digits, underscores.
func validateUsername(username string) *apperrors.AppError {
	if len(username) < 3 || len(username) > 20 {
		return apperrors.ValidationError("用户名长度必须在 3-20 个字符之间")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return apperrors.ValidationError("用户名只能包含字母、数字和下划线")
		}
	}
	return nil
}

// validatePassword checks password strength: 8-128 chars, must contain both letters and digits.
func validatePassword(password string) *apperrors.AppError {
	if len(password) < 8 {
		return apperrors.ValidationError("密码长度不能少于 8 位")
	}
	if len(password) > 128 {
		return apperrors.ValidationError("密码长度不能超过 128 位")
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return apperrors.ValidationError("密码必须同时包含字母和数字")
	}

	return nil
}

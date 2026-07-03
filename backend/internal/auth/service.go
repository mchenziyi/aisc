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

const (
	RefreshTokenBytes = 32 // 64 hex chars
)

type Service struct {
	repo          *Repository
	jwtSecret     []byte
	jwtExpiration time.Duration // access token lifetime
	tokenExpiry   time.Duration // refresh token lifetime
}

func NewService(repo *Repository, jwtSecret string, jwtExpiration time.Duration, tokenExpiry time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
		tokenExpiry:   tokenExpiry,
	}
}

// Register creates a new user account and returns UserWithToken (注册即登录).
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*model.UserWithToken, *apperrors.AppError) {
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

	// Generate tokens (注册即登录)
	tokenResp, appErr := s.generateUserWithToken(ctx, user.ID, user.Username)
	if appErr != nil {
		return nil, appErr
	}

	return tokenResp, nil
}

// Login authenticates a user and returns UserWithToken.
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*model.UserWithToken, *apperrors.AppError) {
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

	// Generate tokens
	tokenResp, appErr := s.generateUserWithToken(ctx, user.ID, user.Username)
	if appErr != nil {
		return nil, appErr
	}

	return tokenResp, nil
}

// RefreshToken validates the refresh token and issues a new access token + refresh token.
func (s *Service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*model.UserWithToken, *apperrors.AppError) {
	// Find user by refresh token hash (we need to scan all users with non-null refresh token)
	// Since we store bcrypt hash, we can't do a direct lookup. Instead, we iterate.
	// A more scalable approach would store a hash we can index, but for simplicity,
	// we'll find the user by checking the refresh token.
	// In practice, you'd use a lookup table. For this exercise, we check all users
	// with a valid refresh_token_expires_at.

	// Find user by refresh token (scan users with active refresh tokens)
	user, err := s.repo.FindByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		log.Printf("internal error: failed to find user by refresh token: %v", err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrRefreshInvalid
	}

	// Check if refresh token is expired
	if user.RefreshTokenExpires != nil && user.RefreshTokenExpires.Before(time.Now()) {
		return nil, apperrors.ErrRefreshInvalid
	}

	// Verify the refresh token against stored hash
	if user.RefreshTokenHash != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(*user.RefreshTokenHash), []byte(req.RefreshToken)); err != nil {
			return nil, apperrors.ErrRefreshInvalid
		}
	} else {
		return nil, apperrors.ErrRefreshInvalid
	}

	// Generate new tokens
	tokenResp, appErr := s.generateUserWithToken(ctx, user.ID, user.Username)
	if appErr != nil {
		return nil, appErr
	}

	return tokenResp, nil
}

// GetMe returns the current user's public info.
func (s *Service) GetMe(ctx context.Context, userID int64) (*model.UserResponse, *apperrors.AppError) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("internal error: failed to find user by ID %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// generateUserWithToken creates JWT + refresh token and returns UserWithToken.
func (s *Service) generateUserWithToken(ctx context.Context, userID int64, username string) (*model.UserWithToken, *apperrors.AppError) {
	now := time.Now()
	expiresIn := s.jwtExpiration

	// Create JWT claims
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Sign JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		log.Printf("internal error: JWT generation failed for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}

	// Generate refresh token (random bytes)
	refreshBytes := make([]byte, RefreshTokenBytes)
	if _, err := rand.Read(refreshBytes); err != nil {
		log.Printf("internal error: failed to generate refresh token for user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}
	refreshToken := hex.EncodeToString(refreshBytes)

	// Store refresh token hash
	s.storeRefreshToken(ctx, userID, refreshToken)

	return &model.UserWithToken{
		Token:     tokenStr,
		ExpiresIn: int(expiresIn.Seconds()),
		ID:        userID,
		Username:  username,
	}, nil
}

// storeRefreshToken hashes and stores the refresh token in the database.
func (s *Service) storeRefreshToken(ctx context.Context, userID int64, refreshToken string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("internal error: failed to hash refresh token for user %d: %v", userID, err)
		return
	}
	hashStr := string(hash)
	expiresAt := time.Now().Add(s.tokenExpiry)

	if err := s.repo.UpdateRefreshToken(ctx, userID, hashStr, expiresAt); err != nil {
		log.Printf("internal error: failed to store refresh token for user %d: %v", userID, err)
	}
}

// validateUsername checks that the username is 3-32 characters and only letters/digits.
func validateUsername(username string) *apperrors.AppError {
	if len(username) < 3 || len(username) > 32 {
		return apperrors.ValidationError("用户名长度必须在 3-32 个字符之间")
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return apperrors.ValidationError("用户名只能包含字母和数字")
		}
	}
	return nil
}

// validatePassword checks password length (8-128).
func validatePassword(password string) *apperrors.AppError {
	if len(password) < 8 {
		return apperrors.ValidationError("密码长度不能少于 8 位")
	}
	if len(password) > 128 {
		return apperrors.ValidationError("密码长度不能超过 128 位")
	}
	return nil
}

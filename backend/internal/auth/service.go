package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
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
	refreshRepo   *RefreshTokenRepo
	jwtSecret     []byte
	jwtExpiration time.Duration // access token lifetime
	tokenExpiry   time.Duration // refresh token lifetime
}

func NewService(repo *Repository, refreshRepo *RefreshTokenRepo, jwtSecret string, jwtExpiration time.Duration, tokenExpiry time.Duration) *Service {
	return &Service{
		repo:          repo,
		refreshRepo:   refreshRepo,
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
	// Find the refresh token record by token hash (indexed lookup via SHA-256)
	record, err := s.refreshRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		log.Printf("refresh token lookup failed: %v", err)
		return nil, apperrors.ErrRefreshInvalid
	}
	if record == nil {
		return nil, apperrors.ErrRefreshInvalid
	}

	// Check if expired
	if record.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrRefreshInvalid
	}

	// Lookup user
	user, err := s.repo.FindByID(ctx, record.UserID)
	if err != nil {
		log.Printf("internal error: failed to find user %d: %v", record.UserID, err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	// Delete the old refresh token (rotation)
	if err := s.refreshRepo.DeleteByUser(ctx, user.ID); err != nil {
		log.Printf("internal error: failed to delete old refresh tokens for user %d: %v", user.ID, err)
		// Continue anyway, the old token will expire
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
		ID:       user.ID,
		Username: user.Username,
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

	// Store refresh token hash in refresh_tokens table (SHA-256 for indexed lookup)
	s.storeRefreshToken(ctx, userID, refreshToken)

	return &model.UserWithToken{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		User: model.UserPublic{
			ID:       userID,
			Username: username,
		},
	}, nil
}

// storeRefreshToken hashes (SHA-256) and stores the refresh token in the refresh_tokens table.
func (s *Service) storeRefreshToken(ctx context.Context, userID int64, refreshToken string) {
	// Use SHA-256 for indexed lookup (unlike bcrypt which can't be indexed)
	h := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(h[:])
	expiresAt := time.Now().Add(s.tokenExpiry)

	if err := s.refreshRepo.Create(ctx, userID, tokenHash, expiresAt); err != nil {
		log.Printf("internal error: failed to store refresh token for user %d: %v", userID, err)
	}
}

// validateUsername checks that the username is 3-20 characters and only letters/digits/underscores.
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

// validatePassword checks password strength: at least 8 chars, must contain letters and digits.
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
		return apperrors.ValidationError("密码必须包含字母和数字")
	}

	return nil
}

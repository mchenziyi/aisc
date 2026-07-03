package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	apperrors "todo-api/internal/errors"
	"todo-api/internal/model"
)

const (
	// RefreshTokenBytes is the byte length of the generated refresh token.
	RefreshTokenBytes = 32 // 64 hex chars
)

type Service struct {
	repo          *Repository
	jwtSecret     []byte
	jwtExpiration time.Duration
	tokenExpiry   time.Duration // refresh token expiry
}

func NewService(repo *Repository, jwtSecret string, jwtExpiration time.Duration, tokenExpiry time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
		tokenExpiry:   tokenExpiry,
	}
}

// Register creates a new user account, logs them in, and returns a UserWithToken.
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*model.UserWithToken, *apperrors.AppError) {
	// Normalize username
	username := strings.ToLower(req.Username)

	// Check if username already exists (fast path before expensive bcrypt)
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
		// Check for unique constraint violation (PostgreSQL error code 23505)
		if code := extractPGErrorCode(err); code == "23505" {
			return nil, apperrors.ErrUsernameTaken
		}
		log.Printf("internal error: failed to create user '%s': %v", username, err)
		return nil, apperrors.ErrInternal
	}

	// Generate tokens
	token, expiresIn, refreshToken, appErr := s.generateTokens(ctx, user.ID)
	if appErr != nil {
		return nil, appErr
	}

	// Store refresh token hash in the background
	s.storeRefreshToken(context.Background(), user.ID, refreshToken)

	return &model.UserWithToken{
		Token:     token,
		ExpiresIn: int(expiresIn.Seconds()),
		ID:        user.ID,
		Username:  user.Username,
	}, nil
}

// Login authenticates a user and returns a UserWithToken.
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
	token, expiresIn, refreshToken, appErr := s.generateTokens(ctx, user.ID)
	if appErr != nil {
		return nil, appErr
	}

	// Store refresh token hash
	s.storeRefreshToken(ctx, user.ID, refreshToken)

	return &model.UserWithToken{
		Token:     token,
		ExpiresIn: int(expiresIn.Seconds()),
		ID:        user.ID,
		Username:  user.Username,
	}, nil
}

// RefreshToken validates the current JWT (via userID) and returns new tokens.
func (s *Service) RefreshToken(ctx context.Context, userID int64) (*model.UserWithToken, *apperrors.AppError) {
	// Verify user still exists
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("internal error: failed to find user %d: %v", userID, err)
		return nil, apperrors.ErrInternal
	}
	if user == nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Generate new tokens
	token, expiresIn, refreshToken, appErr := s.generateTokens(ctx, user.ID)
	if appErr != nil {
		return nil, appErr
	}

	// Store new refresh token (rotate)
	s.storeRefreshToken(ctx, user.ID, refreshToken)

	return &model.UserWithToken{
		Token:     token,
		ExpiresIn: int(expiresIn.Seconds()),
		ID:        user.ID,
		Username:  user.Username,
	}, nil
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

// GenerateAccessToken creates a JWT access token for the given user ID.
// This is used by the middleware for parsing, and by the service for generation.
func (s *Service) GenerateAccessToken(userID int64) (string, time.Duration, error) {
	now := time.Now()
	expiresIn := s.jwtExpiration
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", 0, err
	}
	return token, expiresIn, nil
}

// generateTokens creates a JWT access token and a refresh token.
func (s *Service) generateTokens(ctx context.Context, userID int64) (token string, expiresIn time.Duration, refreshToken string, appErr *apperrors.AppError) {
	// Generate access token
	now := time.Now()
	expiresIn = s.jwtExpiration
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		log.Printf("internal error: JWT generation failed for user %d: %v", userID, err)
		return "", 0, "", apperrors.ErrInternal
	}

	// Generate refresh token (random bytes)
	refreshBytes := make([]byte, RefreshTokenBytes)
	if _, err := rand.Read(refreshBytes); err != nil {
		log.Printf("internal error: failed to generate refresh token for user %d: %v", userID, err)
		return "", 0, "", apperrors.ErrInternal
	}
	refreshToken = hex.EncodeToString(refreshBytes)

	return token, expiresIn, refreshToken, nil
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

// extractPGErrorCode extracts the PostgreSQL error code from an error.
func extractPGErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var pgErr *pgconn.PgError
	if e, ok := err.(*pgconn.PgError); ok {
		return e.Code
	}
	_ = pgErr
	return ""
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if e, ok := err.(*pgconn.PgError); ok {
		return e.Code == "23505"
	}
	_ = pgErr
	return false
}

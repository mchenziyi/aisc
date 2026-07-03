package auth

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	apperrors "todo-api/internal/errors"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	// At least one letter and one digit
	hasLetter = regexp.MustCompile(`[a-zA-Z]`)
	hasDigit  = regexp.MustCompile(`\d`)
)

type Service struct {
	repo          *Repository
	jwtSecret     []byte
	jwtExpiration time.Duration
}

func NewService(repo *Repository, jwtSecret string, jwtExpiration time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiration: jwtExpiration,
	}
}

// Register creates a new user account and returns a JWT token.
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, *apperrors.AppError) {
	// Validate username
	if !usernameRegex.MatchString(req.Username) {
		return nil, apperrors.NewValidationError(
			"username must be 3-20 characters, allowing letters, digits and underscores",
		)
	}

	// Validate password
	if err := validatePassword(req.Password); err != nil {
		return nil, apperrors.NewValidationError(err.Error())
	}

	// Normalize username to lowercase
	username := strings.ToLower(req.Username)

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, username, string(hash))
	if err != nil {
		// Check for unique constraint violation (PostgreSQL error code 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.NewConflictError(apperrors.ErrorCodeUsernameTaken, "username already exists")
		}
		return nil, apperrors.NewInternalError()
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	return &RegisterResponse{
		Token: token,
		User: UserPublic{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, *apperrors.AppError) {
	username := strings.ToLower(req.Username)

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}
	if user == nil {
		return nil, apperrors.NewUnauthorizedError("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid username or password")
	}

	// Generate JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}

	return &LoginResponse{
		Token: token,
		User: UserPublic{
			ID:       user.ID,
			Username: user.Username,
		},
	}, nil
}

// GetMe returns the current user's public info based on user ID.
func (s *Service) GetMe(ctx context.Context, userID int64) (*UserPublic, *apperrors.AppError) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError()
	}
	if user == nil {
		return nil, apperrors.NewNotFoundError("user not found")
	}
	return &UserPublic{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// generateToken creates a JWT token for the given user ID.
func (s *Service) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.jwtExpiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// validatePassword checks password strength rules.
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	letterOk := hasLetter.MatchString(password)
	digitOk := hasDigit.MatchString(password)
	if !letterOk && !digitOk {
		return errors.New("password must contain at least one letter and one digit")
	}
	if !letterOk {
		return errors.New("password must contain at least one letter")
	}
	if !digitOk {
		return errors.New("password must contain at least one digit")
	}
	return nil
}

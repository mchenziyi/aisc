package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUsernameTaken = errors.New("username already exists")
	ErrUserNotFound  = errors.New("user not found")
)

// Repository handles database operations for users.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new auth repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateUser inserts a new user into the database.
func (r *Repository) CreateUser(ctx context.Context, username, passwordHash string) (*User, error) {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)
		RETURNING id, username, created_at, updated_at`

	user := &User{}
	err := r.pool.QueryRow(ctx, query, username, passwordHash).Scan(
		&user.ID, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username (case-insensitive via LOWER).
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at
		FROM users WHERE LOWER(username) = LOWER($1)`

	user := &User{}
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *Repository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at
		FROM users WHERE id = $1`

	user := &User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

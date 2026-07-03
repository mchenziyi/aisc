package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateUser inserts a new user and returns the created user.
func (r *Repository) CreateUser(ctx context.Context, username, passwordHash string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, password) VALUES ($1, $2)
		 RETURNING id, username, password, created_at, updated_at`,
		username, passwordHash,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username (case-insensitive).
func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password, refresh_token_hash, refresh_token_expires_at, created_at, updated_at
		 FROM users WHERE LOWER(username) = LOWER($1)`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.RefreshTokenHash, &user.RefreshTokenExpires, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by their ID.
func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, password, refresh_token_hash, refresh_token_expires_at, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.RefreshTokenHash, &user.RefreshTokenExpires, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateRefreshToken updates the refresh token hash and expiry for a user.
func (r *Repository) UpdateRefreshToken(ctx context.Context, userID int64, hash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET refresh_token_hash = $1, refresh_token_expires_at = $2 WHERE id = $3`,
		hash, expiresAt, userID,
	)
	return err
}

// ClearRefreshToken clears the refresh token fields for a user.
func (r *Repository) ClearRefreshToken(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET refresh_token_hash = NULL, refresh_token_expires_at = NULL WHERE id = $1`,
		userID,
	)
	return err
}

package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

// FindByRefreshToken finds a user whose stored refresh token hash matches the given token.
// Since we use bcrypt, we scan users with non-null refresh_token_hash and compare.
func (r *Repository) FindByRefreshToken(ctx context.Context, refreshToken string) (*User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, username, password, refresh_token_hash, refresh_token_expires_at, created_at, updated_at
		 FROM users WHERE refresh_token_hash IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.RefreshTokenHash, &user.RefreshTokenExpires, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		if user.RefreshTokenHash != nil {
			if err := bcrypt.CompareHashAndPassword([]byte(*user.RefreshTokenHash), []byte(refreshToken)); err == nil {
				return &user, nil
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateRefreshToken updates the refresh token hash and its expiration for a user.
func (r *Repository) UpdateRefreshToken(ctx context.Context, userID int64, refreshTokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET refresh_token_hash = $1, refresh_token_expires_at = $2 WHERE id = $3`,
		refreshTokenHash, expiresAt, userID,
	)
	return err
}

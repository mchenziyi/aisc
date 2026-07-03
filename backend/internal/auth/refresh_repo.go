package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepo struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{pool: pool}
}

// hashToken returns the SHA-256 hex digest of the token.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Create inserts a new refresh token record.
func (r *RefreshTokenRepo) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

// FindByToken finds a refresh token record by the plain refresh token.
// It hashes the token and looks up the hash (indexed lookup, no full scan).
func (r *RefreshTokenRepo) FindByToken(ctx context.Context, refreshToken string) (*RefreshTokenRecord, error) {
	tokenHash := hashToken(refreshToken)

	var record RefreshTokenRecord
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at
		 FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&record.ID, &record.UserID, &record.TokenHash, &record.ExpiresAt, &record.CreatedAt)
	if err != nil {
		return nil, err // pgx.ErrNoRows will propagate
	}
	return &record, nil
}

// DeleteByUser deletes all refresh tokens for a given user (e.g., on password change or logout).
func (r *RefreshTokenRepo) DeleteByUser(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	return err
}

// DeleteExpired cleans up expired refresh tokens.
func (r *RefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE expires_at < NOW()`,
	)
	return err
}

// RefreshTokenRecord represents a refresh token database record.
type RefreshTokenRecord struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

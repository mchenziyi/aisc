package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations applies pending SQL migration files using version tracking.
// It creates a schema_migrations table to track which migrations have been applied.
func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
	// Ensure migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory %s does not exist", migrationsDir)
	}

	// Create schema_migrations table if not exists
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version     VARCHAR(255) PRIMARY KEY,
			applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Read migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory %s: %w", migrationsDir, err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, fname := range upFiles {
		// Extract version from filename (format: 20250301000001_name.up.sql)
		version := extractVersion(fname)
		if version == "" {
			log.Printf("warning: skipping migration %s: could not extract version", fname)
			continue
		}

		// Check if already applied
		var exists bool
		err := pool.QueryRow(context.Background(),
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", fname, err)
		}
		if exists {
			log.Printf("migration skipped (already applied): %s", fname)
			continue
		}

		// Read and execute migration
		path := filepath.Join(migrationsDir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fname, err)
		}

		sql := string(content)
		if _, err := pool.Exec(context.Background(), sql); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fname, err)
		}

		// Record migration
		if _, err := pool.Exec(context.Background(),
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", fname, err)
		}

		log.Printf("migration applied: %s", fname)
	}

	return nil
}

// extractVersion extracts the version prefix from a migration filename.
// Example: "20250301000001_create_users.up.sql" → "20250301000001"
func extractVersion(filename string) string {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

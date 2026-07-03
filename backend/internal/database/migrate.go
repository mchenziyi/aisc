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

// RunMigrations executes SQL migration files from the migrations directory.
func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
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
		path := filepath.Join(migrationsDir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fname, err)
		}

		sql := string(content)
		if _, err := pool.Exec(context.Background(), sql); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fname, err)
		}
		log.Printf("migration applied: %s", fname)
	}

	return nil
}

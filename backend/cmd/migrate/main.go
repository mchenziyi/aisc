package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mchenziyi/aisc/backend/internal/config"
	"github.com/mchenziyi/aisc/backend/internal/database"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.RunMigrations(pool); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migrations completed successfully")
}

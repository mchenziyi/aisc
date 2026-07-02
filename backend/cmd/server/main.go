package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mchenziyi/aisc/backend/internal/auth"
	"github.com/mchenziyi/aisc/backend/internal/config"
	"github.com/mchenziyi/aisc/backend/internal/database"
	"github.com/mchenziyi/aisc/backend/internal/middleware"
	"github.com/mchenziyi/aisc/backend/internal/todo"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set log level
	var logLevel slog.Level
	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	// Initialize database
	ctx := context.Background()
	slog.Info("connecting to database...")
	pool, err := database.NewPool(ctx, cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connection established")

	// Run migrations
	slog.Info("running database migrations...")
	if err := database.RunMigrations(pool); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("database migrations completed")

	// Initialize repositories
	authRepo := auth.NewRepository(pool)
	todoRepo := todo.NewRepository(pool)

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpiration)
	todoService := todo.NewService(todoRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)

	// Setup Gin router
	router := gin.New()

	// Global middlewares
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(middleware.CORSConfig{
		AllowedOrigins: cfg.CORSAllowedOrigins,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dbHealthy := "healthy"
		if err := pool.Ping(ctx); err != nil {
			dbHealthy = "unhealthy"
			slog.Error("health check: database ping failed", "error", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"database":  dbHealthy,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	const maxBodySize = 1 << 20 // 1 MB
	v1 := router.Group("/api/v1")
	v1.Use(middleware.MaxBodySize(maxBodySize))
	{
		// Auth routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Me)
		}

		// Todo routes (JWT required)
		todoGroup := v1.Group("/todos")
		todoGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			todoGroup.POST("", todoHandler.Create)
			todoGroup.GET("", todoHandler.List)
			todoGroup.GET("/:todo_id", todoHandler.GetByID)
			todoGroup.PATCH("/:todo_id", todoHandler.Update)
			todoGroup.DELETE("/:todo_id", todoHandler.Delete)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}

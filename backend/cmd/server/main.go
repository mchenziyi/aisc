package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"todo-api/internal/auth"
	"todo-api/internal/config"
	"todo-api/internal/database"
	apperrors "todo-api/internal/errors"
	"todo-api/internal/middleware"
	"todo-api/internal/todo"
)

func main() {
	cfg := config.Load()

	// Initialize database connection pool
	pool, err := database.NewPool(cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

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
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(cfg.LogLevel))
	r.Use(middleware.ErrorMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dbHealthy := true
		if err := pool.Ping(ctx); err != nil {
			dbHealthy = false
		}

		if !dbHealthy {
			requestID, _ := c.Get("request_id")
			rid, _ := requestID.(string)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":       503,
				"error_code": apperrors.ErrorCodeInternal,
				"message":    "database is unhealthy",
				"request_id": rid,
				"details":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"database":  "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no JWT required)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Protected routes (JWT required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Auth routes that require authentication
			authGroup := protected.Group("/auth")
			{
				authGroup.GET("/me", authHandler.Me)
			}

			// Todo routes
			todoGroup := protected.Group("/todos")
			{
				todoGroup.POST("", todoHandler.CreateTodo)
				todoGroup.GET("", todoHandler.ListTodos)
				todoGroup.GET("/:todo_id", todoHandler.GetTodo)
				todoGroup.PATCH("/:todo_id", todoHandler.UpdateTodo)
				todoGroup.DELETE("/:todo_id", todoHandler.DeleteTodo)
			}
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

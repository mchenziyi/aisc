package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-api/internal/auth"
	"todo-api/internal/config"
	apperrors "todo-api/internal/errors"
	"todo-api/internal/middleware"
	"todo-api/internal/model"
	"todo-api/internal/todo"
)

// Setup configures all routes and returns a Gin engine.
func Setup(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	// Initialize repositories
	authRepo := auth.NewRepository(pool)
	todoRepo := todo.NewRepository(pool)

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTExpiration, cfg.TokenExpiry)
	todoService := todo.NewService(todoRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	todoHandler := todo.NewHandler(todoService)

	// Setup Gin router
	r := gin.New()

	// Global middleware
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware(cfg.LogLevel))
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Health check (no auth required) - GET /health
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, model.SuccessResponse{
				Status:   "error",
				Database: "disconnected",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Status:    "ok",
			Database:  "connected",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})

	// ── Public routes (no JWT required) ──────────────────

	// POST /api/v1/auth/register — Register
	r.POST("/api/v1/auth/register", authHandler.Register)

	// POST /api/v1/auth/login — Login
	r.POST("/api/v1/auth/login", authHandler.Login)

	// POST /api/v1/auth/refresh — Refresh Token
	r.POST("/api/v1/auth/refresh", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.RefreshToken)

	// ── Protected routes (JWT required) ─────────────────
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// GET /api/v1/auth/me — Get current user info
		protected.GET("/auth/me", authHandler.GetCurrentUser)

		// Todos CRUD
		protected.POST("/todos", todoHandler.CreateTodo)       // POST   /api/v1/todos
		protected.GET("/todos", todoHandler.ListTodos)         // GET    /api/v1/todos
		protected.GET("/todos/:id", todoHandler.GetTodo)       // GET    /api/v1/todos/:id
		protected.PATCH("/todos/:id", todoHandler.PatchTodo)   // PATCH  /api/v1/todos/:id
		protected.DELETE("/todos/:id", todoHandler.DeleteTodo) // DELETE /api/v1/todos/:id
	}

	// Catch-all for undefined routes (return 404)
	r.NoRoute(func(c *gin.Context) {
		requestID := c.GetString("request_id")
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			ErrorCode: apperrors.CodeNotFound,
			Message:   "接口不存在",
			RequestID: requestID,
		})
	})

	return r
}

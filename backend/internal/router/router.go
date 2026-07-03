package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-api/internal/auth"
	"todo-api/internal/config"
	apperrors "todo-api/internal/errors"
	"todo-api/internal/middleware"
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
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware(cfg.LogLevel))
	r.Use(middleware.ErrorMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Health check (no auth required) - GET /health
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			log.Printf("health check: database ping failed: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes — matching frozen API spec
	v1 := r.Group("/v1")
	{
		// ── Public routes (no JWT required) ──────────────────

		// POST /v1/users — Register
		v1.POST("/users", authHandler.Register)

		// POST /v1/auth/login — Login
		v1.POST("/auth/login", authHandler.Login)

		// POST /v1/auth/refresh — Refresh Token
		v1.POST("/auth/refresh", authHandler.Refresh)

		// ── Protected routes (JWT required) ─────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// GET /v1/users/me — Get current user info
			protected.GET("/users/me", authHandler.GetCurrentUser)

			// Todos CRUD
			protected.POST("/todos", todoHandler.CreateTodo)       // POST /v1/todos
			protected.GET("/todos", todoHandler.ListTodos)         // GET  /v1/todos
			protected.GET("/todos/:id", todoHandler.GetTodo)       // GET  /v1/todos/:id
			protected.PUT("/todos/:id", todoHandler.UpdateTodo)    // PUT  /v1/todos/:id
			protected.PATCH("/todos/:id", todoHandler.CompleteTodo) // PATCH /v1/todos/:id
			protected.DELETE("/todos/:id", todoHandler.DeleteTodo)  // DELETE /v1/todos/:id
		}
	}

	// Catch-all for undefined routes (return 404)
	r.NoRoute(func(c *gin.Context) {
		apperrors.RespondError(c, apperrors.NewAppError(apperrors.CodeNotFound, http.StatusNotFound, "接口不存在"))
	})

	return r
}

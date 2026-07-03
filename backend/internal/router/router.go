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
	"todo-api/internal/model"
	"todo-api/internal/todo"
)

// Setup configures all routes and returns a Gin engine.
func Setup(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	// Initialize repositories
	authRepo := auth.NewRepository(pool)
	refreshRepo := auth.NewRefreshTokenRepo(pool)
	todoRepo := todo.NewRepository(pool)

	// Initialize services
	authService := auth.NewService(authRepo, refreshRepo, cfg.JWTSecret, cfg.JWTExpiration, cfg.TokenExpiry)
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

		dbStatus := "ok"
		if err := pool.Ping(ctx); err != nil {
			log.Printf("health check: database ping failed: %v", err)
			dbStatus = "error"
		}

		c.JSON(http.StatusOK, model.HealthResponse{
			Status:    "ok",
			Database:  dbStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1 routes — matching frozen API spec
	apiV1 := r.Group("/api/v1")
	{
		// ── Public routes (no JWT required) ──────────────────

		// POST /api/v1/auth/register — Register
		apiV1.POST("/auth/register", authHandler.Register)

		// POST /api/v1/auth/login — Login
		apiV1.POST("/auth/login", authHandler.Login)

		// POST /api/v1/auth/refresh — Refresh Token
		apiV1.POST("/auth/refresh", authHandler.Refresh)

		// ── Protected routes (JWT required) ─────────────────
		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// GET /api/v1/auth/me — Get current user info
			protected.GET("/auth/me", authHandler.GetCurrentUser)

			// Todos CRUD
			protected.POST("/todos", todoHandler.CreateTodo)        // POST   /api/v1/todos
			protected.GET("/todos", todoHandler.ListTodos)          // GET    /api/v1/todos
			protected.GET("/todos/:id", todoHandler.GetTodo)        // GET    /api/v1/todos/:id
			protected.PATCH("/todos/:id", todoHandler.PatchTodo)    // PATCH  /api/v1/todos/:id (update with version)
			protected.DELETE("/todos/:id", todoHandler.DeleteTodo)  // DELETE /api/v1/todos/:id (with version)
		}
	}

	// Catch-all for undefined routes (return 404)
	r.NoRoute(func(c *gin.Context) {
		apperrors.RespondError(c, apperrors.NewAppError(apperrors.CodeNotFound, http.StatusNotFound, "接口不存在"))
	})

	return r
}

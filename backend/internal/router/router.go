package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-api/internal/auth"
	"todo-api/internal/config"
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
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(cfg.LogLevel))
	r.Use(middleware.ErrorMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Health check (no auth required)
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dbHealthy := true
		if err := pool.Ping(ctx); err != nil {
			dbHealthy = false
		}

		if !dbHealthy {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Public routes (no JWT required)
		v1.POST("/users", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		// Protected routes (JWT required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Token refresh (requires valid JWT)
			protected.POST("/auth/refresh", authHandler.RefreshToken)

			// Users
			protected.GET("/users/me", authHandler.GetCurrentUser)

			// Todos
			protected.POST("/todos", todoHandler.CreateTodo)
			protected.GET("/todos", todoHandler.ListTodos)
			protected.GET("/todos/:id", todoHandler.GetTodo)
			protected.PUT("/todos/:id", todoHandler.UpdateTodo)
			protected.PATCH("/todos/:id", todoHandler.CompleteTodo)
			protected.DELETE("/todos/:id", todoHandler.DeleteTodo)
		}
	}

	return r
}

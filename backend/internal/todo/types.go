package todo

import (
	"time"
)

// ─── Request DTOs ─────────────────────────────────────────────

// CreateTodoRequest represents the request body for creating a todo.
// Title: 1-255 chars. Description: max 1000 chars. DueDate: ISO 8601 date.
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
	DueDate     *string `json:"due_date,omitempty"` // YYYY-MM-DD
}

// UpdateTodoRequest represents the request body for updating a todo (PATCH).
// All fields are optional; version is required for optimistic locking.
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
	DueDate     *string `json:"due_date,omitempty"` // YYYY-MM-DD
	Completed   *bool   `json:"completed,omitempty"`
	Version     int     `json:"version" binding:"required"`
}

// ─── Query DTOs ──────────────────────────────────────────────

// ListQuery represents query parameters for listing todos.
type ListQuery struct {
	Page     int
	PageSize int
	Status   string // "all", "pending", "completed"
}

// ─── Database Model ───────────────────────────────────────────

// Todo represents a todo record from the database.
type Todo struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
	Version     int        `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

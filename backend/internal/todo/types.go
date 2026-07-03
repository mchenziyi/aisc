package todo

import (
	"time"
)

// ─── Request DTOs ─────────────────────────────────────────────

// CreateTodoRequest represents the request body for creating a todo.
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
	DueDate     *string `json:"due_date,omitempty"` // ISO 8601 format
}

// PatchTodoRequest represents the request body for partially updating a todo.
// All fields are optional. Supports setting completed status.
type PatchTodoRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
	DueDate     *string `json:"due_date,omitempty"` // ISO 8601 format
	Completed   *bool   `json:"completed,omitempty"`
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
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	Version     int        `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

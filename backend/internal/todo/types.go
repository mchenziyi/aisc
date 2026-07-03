package todo

import "time"

// ─── Request DTOs ─────────────────────────────────────────────

// CreateTodoRequest represents the request body for creating a todo.
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=200"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
	DueDate     *string `json:"due_date,omitempty"` // ISO 8601 format, must be future
}

// UpdateTodoRequest represents the request body for updating a todo.
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
	DueDate     *string `json:"due_date,omitempty"` // ISO 8601, nil = no change, "" = clear
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
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

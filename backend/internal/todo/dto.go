package todo

import "time"

// CreateTodoRequest represents the request body for creating a todo.
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"` // format: YYYY-MM-DD
}

// UpdateTodoRequest represents the request body for updating a todo.
type UpdateTodoRequest struct {
	Version     int64   `json:"version" binding:"required"`
	Title       *string `json:"title"`
	Description *string `json:"description"` // empty string ("") clears the field (sets to NULL)
	DueDate     *string `json:"due_date"`    // empty string ("") clears the field (sets to NULL), format: YYYY-MM-DD
	Completed   *bool   `json:"completed"`
}

// DeleteTodoRequest represents the request body for deleting a todo.
type DeleteTodoRequest struct {
	Version int64 `json:"version" binding:"required"`
}

// TodoResponse represents the full todo object returned by the API.
type TodoResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *string    `json:"due_date"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Version     int64      `json:"version"`
	UserID      int64      `json:"user_id"`
}

// TodoListResponse represents the paginated list response.
type TodoListResponse struct {
	Items      []TodoResponse `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

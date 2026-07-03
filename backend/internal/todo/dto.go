package todo

import (
	"encoding/json"
	"time"
)

// ─── NullableString ──────────────────────────────────────────

// NullableString 区分三种 JSON 状态：未提供、null、有值。
// - JSON 字段缺失 → IsSet=false, IsNull=false
// - JSON "field": null → IsSet=true, IsNull=true
// - JSON "field": "hello" → IsSet=true, Value="hello"
type NullableString struct {
	Value  string
	IsSet  bool
	IsNull bool
}

func (n *NullableString) UnmarshalJSON(data []byte) error {
	n.IsSet = true
	if string(data) == "null" {
		n.IsNull = true
		return nil
	}
	return json.Unmarshal(data, &n.Value)
}

// ─── DTOs ────────────────────────────────────────────────────

// CreateTodoRequest represents the request body for creating a todo.
type CreateTodoRequest struct {
	Title       string         `json:"title" binding:"required"`
	Description NullableString `json:"description"`
	DueDate     NullableString `json:"due_date"`
}

// UpdateTodoRequest represents the request body for updating a todo.
type UpdateTodoRequest struct {
	Version     int64          `json:"version" binding:"required"`
	Title       *string        `json:"title"`
	Description NullableString `json:"description"` // null → 清空；未提供 → 不修改
	DueDate     NullableString `json:"due_date"`    // null → 清空；未提供 → 不修改
	Completed   *bool          `json:"completed"`
}

// DeleteTodoRequest represents the request body for deleting a todo.
// Version is optional in body — handler also accepts it from query param.
type DeleteTodoRequest struct {
	Version int64 `json:"version"`
}

// TodoResponse represents the full todo object returned by the API.
type TodoResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	DueDate     *string   `json:"due_date"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int64     `json:"version"`
	UserID      int64     `json:"user_id"`
}

// TodoListResponse represents the paginated list response.
type TodoListResponse struct {
	Items      []TodoResponse `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

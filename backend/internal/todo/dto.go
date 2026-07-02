package todo

import "encoding/json"

// CreateTodoRequest represents the request to create a new todo.
type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"` // format: YYYY-MM-DD
}

// OptionalString represents a nullable string field that distinguishes
// between "not provided", "set to null", and "set to a value".
type OptionalString struct {
	Set   bool
	Null  bool
	Value string
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Null = true
		return nil
	}
	return json.Unmarshal(data, &o.Value)
}

// UpdateTodoRequest represents the request to update a todo.
type UpdateTodoRequest struct {
	Version     int64           `json:"version" binding:"required,min=1"`
	Title       *string         `json:"title,omitempty"`
	Description *OptionalString `json:"description,omitempty"`
	DueDate     *OptionalString `json:"due_date,omitempty"`
	Completed   *bool           `json:"completed,omitempty"`
}

// Todo represents a todo item in the database.
type Todo struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
	Completed   bool    `json:"completed"`
	Version     int64   `json:"version"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// TodoListResponse represents the paginated list response.
type TodoListResponse struct {
	Items      []Todo `json:"items"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// ListParams holds pagination parameters.
type ListParams struct {
	Page     int
	PageSize int
}

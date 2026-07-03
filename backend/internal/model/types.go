package model

// ResponseEnvelope is the unified API response wrapper.
type ResponseEnvelope struct {
	Code    int          `json:"code"`              // 业务码，0 为成功
	Message string       `json:"message"`           // 提示信息
	Data    interface{}  `json:"data,omitempty"`    // 成功时携带数据
	Errors  []FieldError `json:"errors,omitempty"` // 字段级错误列表（可选）
}

// FieldError represents a field-level validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewSuccessResponse creates a success response envelope.
func NewSuccessResponse(data interface{}) *ResponseEnvelope {
	return &ResponseEnvelope{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

// NewErrorResponse creates an error response envelope with optional field errors.
func NewErrorResponse(code int, message string, errors []FieldError) *ResponseEnvelope {
	resp := &ResponseEnvelope{
		Code:    code,
		Message: message,
	}
	if len(errors) > 0 {
		resp.Errors = errors
	}
	return resp
}

// UserResponse represents the public user info (no password).
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// UserWithToken represents the auth response with JWT token and user info.
type UserWithToken struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	ID        int64  `json:"id"`
	Username  string `json:"username"`
}

// TodoResponse represents a complete todo item as returned by the API.
type TodoResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	IsCompleted bool    `json:"is_completed"`
	CompletedAt *string `json:"completed_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	UserID      int64   `json:"user_id"`
}

// PaginatedTodos represents a paginated list of todos.
type PaginatedTodos struct {
	Items      []TodoResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

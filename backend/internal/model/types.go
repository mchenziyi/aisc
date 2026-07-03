package model

// FieldError represents a field-level validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// SuccessResponse is the unified API success response format (used only for health check).
type SuccessResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

// ErrorResponse is the unified API error response format.
type ErrorResponse struct {
	ErrorCode string       `json:"error_code"`          // string error code, e.g. "VALIDATION_ERROR"
	Message   string       `json:"message"`             // human-readable error description
	RequestID string       `json:"request_id"`          // unique request identifier
	Details   []FieldError `json:"details,omitempty"`   // field-level error details (optional)
}

// UserPublic represents public user information (id, username, created_at).
type UserPublic struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// UserInfo represents user info embedded in auth responses.
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// AuthResponse is returned for register and login endpoints.
type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// TodoResponse represents a todo item as returned by the API.
type TodoResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	UserID      int64   `json:"user_id"`
	Version     int     `json:"version"`
}

// PaginatedTodos is the paginated list response for todos.
type PaginatedTodos struct {
	Items      []TodoResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// String error codes matching the API spec.
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeTokenExpired   = "TOKEN_EXPIRED"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeVersionConflict = "VERSION_CONFLICT"
	ErrCodeInternal       = "INTERNAL_ERROR"
)

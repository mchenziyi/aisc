package model

// ─── Response Envelope (REMOVED — 直接返回资源对象) ──────────────

// ─── Error Response ──────────────────────────────────────────

// FieldError represents a field-level validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse is the unified error response format.
// Matches the API spec: { error_code, message, request_id, details? }.
type ErrorResponse struct {
	ErrorCode string       `json:"error_code"`
	Message   string       `json:"message"`
	RequestID string       `json:"request_id"`
	Details   []FieldError `json:"details,omitempty"`
}

// ─── User Schemas ────────────────────────────────────────────

// UserPublic is the minimal public user info.
type UserPublic struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// UserWithToken is returned for register, login, and refresh endpoints.
// Matches the API spec: { token, user: { id, username }, refresh_token }.
type UserWithToken struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
	User         UserPublic `json:"user"`
}

// UserResponse represents public user information (no password).
// Used by GET /auth/me — only id and username.
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// ─── Todo Schemas ────────────────────────────────────────────

// TodoResponse represents a todo item as returned by the API.
// Matches the API spec exactly.
type TodoResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`    // nil → null in JSON
	DueDate     *string `json:"due_date"`       // nil → null in JSON (YYYY-MM-DD)
	Completed   bool    `json:"completed"`
	Version     int     `json:"version"`
	CreatedAt   string  `json:"created_at"`     // ISO 8601
	UpdatedAt   string  `json:"updated_at"`     // ISO 8601
	UserID      int64   `json:"user_id"`
}

// PaginatedTodos is the paginated list response for todos.
// Matches the API spec: { items, total, page, page_size, total_pages }.
type PaginatedTodos struct {
	Items      []TodoResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

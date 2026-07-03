package model

// ─── Response Envelope ────────────────────────────────────────

// ResponseEnvelope is the unified success response format.
// code: 0 for success, non-zero for business errors.
// data: the payload for success responses.
// message: human-readable message ("ok" for success).
type ResponseEnvelope struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// ─── Error Response ──────────────────────────────────────────

// FieldError represents a field-level validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse is the unified error response format.
// Matches the API spec: { code, message, errors? }.
// code is a business error code (e.g. 1001, 2001), NOT HTTP status.
type ErrorResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors,omitempty"`
}

// ─── User Schemas ────────────────────────────────────────────

// UserWithToken is returned for register, login, and refresh endpoints.
// Matches the API spec: { token, expires_in, id, username }.
type UserWithToken struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	ID        int64  `json:"id"`
	Username  string `json:"username"`
}

// UserResponse represents public user information (no password).
// Used by GET /v1/users/me.
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// ─── Todo Schemas ────────────────────────────────────────────

// TodoResponse represents a todo item as returned by the API.
// Matches the API spec exactly.
type TodoResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`    // nil → null in JSON
	DueDate     *string `json:"due_date"`       // nil → null in JSON (ISO 8601)
	IsCompleted bool    `json:"is_completed"`
	CompletedAt *string `json:"completed_at"`   // nil → null in JSON (ISO 8601)
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

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
//
// ⚠️ 空字符串等价于 null 的行为说明：
// 在 Create 和 Update 两个接口中，空字符串 "" 与 null 同等对待，
// 均视为"清除/置空"操作：
//   - Create 时：description 或 due_date 传 "" 表示不设置该字段（等价于不传或传 null）
//   - Update 时：description 或 due_date 传 "" 表示清空现有值（等价于传 null）
//
// 此行为是设计使然，因为 HTTP JSON 语义中空字符串对于可选字段
// 没有明确的业务含义，统一作为"置空"处理可以减少客户端的判断逻辑。
//
// 如果客户端不希望改变某个字段的值，应在请求体中省略该字段（不传），
// 而非传入空字符串或 null。
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
	Title       string         `json:"title" binding:"required,min=1,max=255"`
	Description NullableString `json:"description"`
	DueDate     NullableString `json:"due_date"`
}

// UpdateTodoRequest represents the request body for updating a todo.
type UpdateTodoRequest struct {
	Version     int64          `json:"version" binding:"required,min=1"`
	Title       *string        `json:"title"`
	Description NullableString `json:"description"` // null → 清空；未提供 → 不修改
	DueDate     NullableString `json:"due_date"`    // null → 清空；未提供 → 不修改
	Completed   *bool          `json:"completed"`
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

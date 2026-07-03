package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "todo-api/internal/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateTodo handles POST /api/v1/todos
func (h *Handler) CreateTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		c.Abort()
		return
	}

	resp, appErr := h.service.Create(c.Request.Context(), userID, &req)
	if appErr != nil {
		c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListTodos handles GET /api/v1/todos
func (h *Handler) ListTodos(c *gin.Context) {
	userID := c.GetInt64("user_id")

	page, err := parseQueryInt(c, "page", 1)
	if err != nil {
		c.Error(apperrors.NewValidationError("page must be a valid integer"))
		c.Abort()
		return
	}
	if page < 1 {
		c.Error(apperrors.NewValidationError("page must be >= 1"))
		c.Abort()
		return
	}

	pageSize, err := parseQueryInt(c, "page_size", 20)
	if err != nil {
		c.Error(apperrors.NewValidationError("page_size must be a valid integer"))
		c.Abort()
		return
	}
	if pageSize < 1 {
		c.Error(apperrors.NewValidationError("page_size must be >= 1"))
		c.Abort()
		return
	}
	if pageSize > 100 {
		c.Error(apperrors.NewValidationError("page_size must not exceed 100"))
		c.Abort()
		return
	}

	resp, appErr := h.service.List(c.Request.Context(), userID, page, pageSize)
	if appErr != nil {
		c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTodo handles GET /api/v1/todos/:todo_id
func (h *Handler) GetTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("invalid todo_id format"))
		c.Abort()
		return
	}
	if todoID < 1 {
		c.Error(apperrors.NewValidationError("todo_id must be a positive integer"))
		c.Abort()
		return
	}

	resp, appErr := h.service.GetByID(c.Request.Context(), userID, todoID)
	if appErr != nil {
		c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateTodo handles PATCH /api/v1/todos/:todo_id
func (h *Handler) UpdateTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("invalid todo_id format"))
		c.Abort()
		return
	}
	if todoID < 1 {
		c.Error(apperrors.NewValidationError("todo_id must be a positive integer"))
		c.Abort()
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperrors.NewValidationErrorFromBinding(err))
		c.Abort()
		return
	}

	// Validate version >= 1
	if req.Version < 1 {
		c.Error(apperrors.NewValidationError("version must be >= 1"))
		c.Abort()
		return
	}

	// Validate that at least one field (other than version) is provided
	if req.Title == nil && !req.Description.IsSet && !req.DueDate.IsSet && req.Completed == nil {
		c.Error(apperrors.NewValidationError("at least one field to update (other than version) must be provided"))
		c.Abort()
		return
	}

	resp, appErr := h.service.Update(c.Request.Context(), userID, todoID, &req)
	if appErr != nil {
		c.Error(appErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTodo handles DELETE /api/v1/todos/:todo_id
// Version is required as a query parameter for optimistic locking.
func (h *Handler) DeleteTodo(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("invalid todo_id format"))
		c.Abort()
		return
	}
	if todoID < 1 {
		c.Error(apperrors.NewValidationError("todo_id must be a positive integer"))
		c.Abort()
		return
	}

	versionStr := c.Query("version")
	if versionStr == "" {
		c.Error(apperrors.NewValidationError("version is required and must be >= 1"))
		c.Abort()
		return
	}
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		c.Error(apperrors.NewValidationError("version must be a valid integer"))
		c.Abort()
		return
	}
	if version < 1 {
		c.Error(apperrors.NewValidationError("version must be >= 1"))
		c.Abort()
		return
	}

	appErr := h.service.Delete(c.Request.Context(), userID, todoID, version)
	if appErr != nil {
		c.Error(appErr)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

// parseQueryInt parses an integer query parameter with a default value.
func parseQueryInt(c *gin.Context, key string, defaultVal int) (int, error) {
	val := c.DefaultQuery(key, "")
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}

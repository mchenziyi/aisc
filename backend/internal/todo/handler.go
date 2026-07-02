package todo

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler handles todo HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new todo handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /todos
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	todo, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		if isValidationError(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// List handles GET /todos
func (h *Handler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := h.service.List(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		if isValidationError(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetByID handles GET /todos/:todo_id
func (h *Handler) GetByID(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid todo_id format",
		})
		return
	}

	todo, err := h.service.GetByID(c.Request.Context(), userID, todoID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "todo not found",
			})
			return
		}
		if errors.Is(err, ErrForbidden) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "forbidden: todo belongs to another user",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Update handles PATCH /todos/:todo_id
func (h *Handler) Update(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid todo_id format",
		})
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	todo, err := h.service.Update(c.Request.Context(), userID, todoID, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "todo not found",
			})
			return
		}
		if errors.Is(err, ErrForbidden) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "forbidden: todo belongs to another user",
			})
			return
		}
		if errors.Is(err, ErrVersionConflict) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "resource conflict due to version mismatch",
			})
			return
		}
		// Validation errors
		if isValidationError(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}
		// Unexpected server errors
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Delete handles DELETE /todos/:todo_id
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")

	todoID, err := strconv.ParseInt(c.Param("todo_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid todo_id format",
		})
		return
	}

	versionStr := c.Query("version")
	if versionStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "version query parameter is required",
		})
		return
	}
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil || version < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "version must be a positive integer",
		})
		return
	}

	err = h.service.Delete(c.Request.Context(), userID, todoID, version)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "todo not found",
			})
			return
		}
		if errors.Is(err, ErrForbidden) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "forbidden: todo belongs to another user",
			})
			return
		}
		if errors.Is(err, ErrVersionConflict) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "resource conflict due to version mismatch",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// isValidationError checks if an error is a client-side validation error.
func isValidationError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "must be") ||
		strings.Contains(msg, "is required") ||
		strings.Contains(msg, "cannot be") ||
		strings.Contains(msg, "at least one field")
}

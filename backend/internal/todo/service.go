package todo

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

var (
	dateRegex       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	ErrNotFound     = errors.New("todo not found")
	ErrForbidden    = errors.New("forbidden: todo belongs to another user")
	ErrVersionConflict = errors.New("resource conflict due to version mismatch")
)

// Service handles todo business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new todo service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new todo for the user.
func (s *Service) Create(ctx context.Context, userID int64, req CreateTodoRequest) (*Todo, error) {
	// Validate title
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if len(title) > 255 {
		return nil, fmt.Errorf("title must be at most 255 characters")
	}

	// Validate description length if provided
	if req.Description != nil && len(*req.Description) > 1000 {
		return nil, fmt.Errorf("description must be at most 1000 characters")
	}

	// Validate due_date format if provided
	if req.DueDate != nil && *req.DueDate != "" {
		if !dateRegex.MatchString(*req.DueDate) {
			return nil, fmt.Errorf("due_date must be in YYYY-MM-DD format")
		}
		// Also verify it's a valid date
		if _, err := time.Parse("2006-01-02", *req.DueDate); err != nil {
			return nil, fmt.Errorf("due_date must be a valid date in YYYY-MM-DD format")
		}
	}

	todo := &Todo{
		UserID:      userID,
		Title:       title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	result, err := s.repo.Create(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}
	return result, nil
}

// List retrieves a paginated list of todos for the user.
func (s *Service) List(ctx context.Context, userID int64, page, pageSize int) (*TodoListResponse, error) {
	if page < 1 {
		return nil, fmt.Errorf("page must be >= 1")
	}
	if pageSize < 1 {
		return nil, fmt.Errorf("page_size must be >= 1")
	}
	if pageSize > 100 {
		return nil, fmt.Errorf("page_size must be <= 100")
	}

	params := ListParams{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := s.repo.List(ctx, userID, params)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	return resp, nil
}

// GetByID retrieves a single todo for the user.
// Returns (nil, ErrNotFound) if not found, or an error with a message for forbidden access.
func (s *Service) GetByID(ctx context.Context, userID, todoID int64) (*Todo, error) {
	// First check if the todo exists at all (without user filter)
	todo, err := s.repo.GetByIDAnyUser(ctx, todoID)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	// Check ownership
	if todo.UserID != userID {
		return nil, ErrForbidden
	}

	return todo, nil
}

// GetByIDOwned retrieves a todo only if it belongs to the user.
// Returns (nil, ErrNotFound) if not found or not owned.
func (s *Service) GetByIDOwned(ctx context.Context, userID, todoID int64) (*Todo, error) {
	todo, err := s.repo.GetByID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// Update updates a todo with optimistic locking.
func (s *Service) Update(ctx context.Context, userID, todoID int64, req UpdateTodoRequest) (*Todo, error) {
	// At least one field (other than version) must be provided
	if req.Title == nil && req.Description == nil && req.DueDate == nil && req.Completed == nil {
		return nil, fmt.Errorf("at least one field to update (other than version) must be provided")
	}

	// Validate title if provided
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
		if len(title) > 255 {
			return nil, fmt.Errorf("title must be at most 255 characters")
		}
		req.Title = &title
	}

	// Validate description length if provided
	if req.Description != nil && !req.Description.Null && req.Description.Value != "" && len(req.Description.Value) > 1000 {
		return nil, fmt.Errorf("description must be at most 1000 characters")
	}

	// Validate due_date if provided
	if req.DueDate != nil && !req.DueDate.Null && req.DueDate.Value != "" {
		if !dateRegex.MatchString(req.DueDate.Value) {
			return nil, fmt.Errorf("due_date must be in YYYY-MM-DD format")
		}
		if _, err := time.Parse("2006-01-02", req.DueDate.Value); err != nil {
			return nil, fmt.Errorf("due_date must be a valid date in YYYY-MM-DD format")
		}
	}

	// First check if the todo exists at all
	existingTodo, err := s.repo.GetByIDAnyUser(ctx, todoID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	// Check ownership
	if existingTodo.UserID != userID {
		return nil, ErrForbidden
	}

	// Check version match
	if existingTodo.Version != req.Version {
		return nil, ErrVersionConflict
	}

	// Idempotency check for completed-only updates
	if req.Title == nil && req.Description == nil && req.DueDate == nil && req.Completed != nil {
		if *req.Completed == existingTodo.Completed {
			// State matches, no update needed (idempotent)
			return existingTodo, nil
		}
		// State differs, just update completed
		result, err := s.repo.UpdateCompletedOnly(ctx, todoID, userID, req.Version, *req.Completed)
		if err != nil {
			if errors.Is(err, ErrVersionConflict) {
				return nil, ErrVersionConflict
			}
			return nil, fmt.Errorf("update completed: %w", err)
		}
		return result, nil
	}

	// General update with multiple fields
	// For nullable fields (description, due_date), we need to handle NULL explicitly
	fields := UpdateFields{}

	if req.Title != nil {
		fields.Title = req.Title
	}

	if req.Description != nil {
		if req.Description.Null {
			fields.Description = &FieldString{SetNull: true}
		} else {
			fields.Description = &FieldString{Value: req.Description.Value}
		}
	}

	if req.DueDate != nil {
		if req.DueDate.Null {
			fields.DueDate = &FieldString{SetNull: true}
		} else {
			fields.DueDate = &FieldString{Value: req.DueDate.Value}
		}
	}

	if req.Completed != nil {
		fields.Completed = req.Completed
	}

	result, err := s.repo.UpdateWithFields(ctx, todoID, userID, req.Version, fields)
	if err != nil {
		if errors.Is(err, ErrVersionConflict) {
			return nil, ErrVersionConflict
		}
		return nil, fmt.Errorf("update todo: %w", err)
	}
	return result, nil
}

// Delete deletes a todo with optimistic locking.
func (s *Service) Delete(ctx context.Context, userID, todoID int64, version int64) error {
	// First check if the todo exists at all
	existingTodo, err := s.repo.GetByIDAnyUser(ctx, todoID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("get todo: %w", err)
	}

	// Check ownership
	if existingTodo.UserID != userID {
		return ErrForbidden
	}

	// Attempt delete with version check
	err = s.repo.Delete(ctx, todoID, userID, version)
	if err != nil {
		if errors.Is(err, ErrVersionConflict) {
			return ErrVersionConflict
		}
		return fmt.Errorf("delete todo: %w", err)
	}
	return nil
}

// validateAndNormalizePage validates pagination parameters.
func validateAndNormalizePage(page, pageSize int) (int, int, error) {
	if page < 1 {
		return 0, 0, fmt.Errorf("page must be >= 1")
	}
	if pageSize < 1 {
		return 0, 0, fmt.Errorf("page_size must be >= 1")
	}
	if pageSize > 100 {
		return 0, 0, fmt.Errorf("page_size must be <= 100")
	}
	return page, pageSize, nil
}

// calculateTotalPages computes total pages.
func calculateTotalPages(total, pageSize int) int {
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

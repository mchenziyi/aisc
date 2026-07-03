package todo

import (
	"context"
	"errors"
	"log"
	"regexp"
	"time"
	"unicode/utf8"

	apperrors "todo-api/internal/errors"
)

var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var errInvalidDateFormat = errors.New("due_date must be in YYYY-MM-DD format")

const (
	maxTitleLen       = 255
	maxDescriptionLen = 1000
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create creates a new todo for the given user.
func (s *Service) Create(ctx context.Context, userID int64, req *CreateTodoRequest) (*TodoResponse, *apperrors.AppError) {
	// Validate title
	if req.Title == "" {
		return nil, apperrors.NewValidationError("title is required")
	}
	if utf8.RuneCountInString(req.Title) > maxTitleLen {
		return nil, apperrors.NewValidationError("title must not exceed 255 characters")
	}

	// Validate description length
	if req.Description.IsSet && !req.Description.IsNull && utf8.RuneCountInString(req.Description.Value) > maxDescriptionLen {
		return nil, apperrors.NewValidationError("description must not exceed 1000 characters")
	}

	// Parse and validate due_date if provided
	var dueDate *time.Time
	if req.DueDate.IsSet && !req.DueDate.IsNull && req.DueDate.Value != "" {
		parsed, err := parseDate(req.DueDate.Value)
		if err != nil {
			return nil, apperrors.NewValidationError(err.Error())
		}
		dueDate = &parsed
	}

	// Build description pointer from NullableString
	// Empty string is treated as null (clear the field)
	var desc *string
	if req.Description.IsSet && !req.Description.IsNull && req.Description.Value != "" {
		desc = &req.Description.Value
	}

	todo := &Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: desc,
		DueDate:     dueDate,
		Completed:   false,
		Version:     1,
	}

	created, err := s.repo.Create(ctx, todo)
	if err != nil {
		log.Printf("internal error: failed to create todo for user %d: %v", userID, err)
		return nil, apperrors.NewInternalError()
	}

	return toResponse(created), nil
}

// GetByID retrieves a single todo by ID, ensuring it belongs to the user.
func (s *Service) GetByID(ctx context.Context, userID, todoID int64) (*TodoResponse, *apperrors.AppError) {
	todo, err := s.repo.FindByIDAndUser(ctx, todoID, userID)
	if err != nil {
		log.Printf("internal error: failed to find todo %d for user %d: %v", todoID, userID, err)
		return nil, apperrors.NewInternalError()
	}
	if todo == nil {
		return nil, apperrors.NewNotFoundError("todo not found")
	}
	return toResponse(todo), nil
}

// List returns a paginated list of todos for the user.
func (s *Service) List(ctx context.Context, userID int64, page, pageSize int) (*TodoListResponse, *apperrors.AppError) {
	todos, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		log.Printf("internal error: failed to list todos for user %d: %v", userID, err)
		return nil, apperrors.NewInternalError()
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	items := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		items = append(items, *toResponse(t))
	}

	return &TodoListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a todo with optimistic locking.
func (s *Service) Update(ctx context.Context, userID, todoID int64, req *UpdateTodoRequest) (*TodoResponse, *apperrors.AppError) {
	// Validate version
	if req.Version < 1 {
		return nil, apperrors.NewValidationError("version must be >= 1")
	}

	// Validate that at least one field (other than version) is provided
	if req.Title == nil && !req.Description.IsSet && !req.DueDate.IsSet && req.Completed == nil {
		return nil, apperrors.NewValidationError("at least one field to update (other than version) must be provided")
	}

	// Validate title if provided
	if req.Title != nil {
		if *req.Title == "" {
			return nil, apperrors.NewValidationError("title cannot be empty")
		}
		if utf8.RuneCountInString(*req.Title) > maxTitleLen {
			return nil, apperrors.NewValidationError("title must not exceed 255 characters")
		}
	}

	// Validate description if provided
	if req.Description.IsSet && !req.Description.IsNull && utf8.RuneCountInString(req.Description.Value) > maxDescriptionLen {
		return nil, apperrors.NewValidationError("description must not exceed 1000 characters")
	}

	// Parse and validate due_date if provided
	var updateDueDate bool
	var dueDateVal *time.Time
	if req.DueDate.IsSet {
		updateDueDate = true
		if !req.DueDate.IsNull && req.DueDate.Value != "" {
			parsed, err := parseDate(req.DueDate.Value)
			if err != nil {
				return nil, apperrors.NewValidationError(err.Error())
			}
			dueDateVal = &parsed
		}
		// If IsNull or Value=="", dueDateVal stays nil → sets to SQL NULL
	}

	// Check idempotent case: if only version + completed is provided
	// and the completed value matches the current state, return without modification.
	if req.Title == nil && !req.Description.IsSet && !req.DueDate.IsSet && req.Completed != nil {
		existing, err := s.repo.FindByIDAndUser(ctx, todoID, userID)
		if err != nil {
			log.Printf("internal error: failed to find todo %d for user %d: %v", todoID, userID, err)
			return nil, apperrors.NewInternalError()
		}
		if existing == nil {
			return nil, apperrors.NewNotFoundError("todo not found")
		}
		if existing.Version != req.Version {
			return nil, apperrors.NewVersionConflictError(existing.Version)
		}
		if existing.Completed == *req.Completed {
			// Idempotent: return current state without modification
			return toResponse(existing), nil
		}
	}

	// Build update fields
	fields := &UpdateFields{
		Title:         req.Title,
		DueDateVal:    dueDateVal,
		UpdateDueDate: updateDueDate,
		Completed:     req.Completed,
	}
	if req.Description.IsSet {
		fields.UpdateDescription = true
		// Empty string is treated as null (clear the field)
		if req.Description.IsNull || req.Description.Value == "" {
			fields.DescriptionVal = nil // clear to NULL
		} else {
			fields.DescriptionVal = &req.Description.Value
		}
	}

	updated, err := s.repo.UpdateVersioned(ctx, todoID, userID, req.Version, fields)
	if err != nil {
		if err == errVersionConflict {
			// Get current version for response
			existing, _ := s.repo.FindByIDAndUser(ctx, todoID, userID)
			if existing != nil {
				return nil, apperrors.NewVersionConflictError(existing.Version)
			}
			return nil, apperrors.NewVersionConflictError(0)
		}
		log.Printf("internal error: failed to update todo %d for user %d: %v", todoID, userID, err)
		return nil, apperrors.NewInternalError()
	}
	if updated == nil {
		log.Printf("internal error: update returned nil for todo %d user %d without error", todoID, userID)
		return nil, apperrors.NewNotFoundError("todo not found")
	}

	return toResponse(updated), nil
}

// Delete deletes a todo with optimistic locking.
func (s *Service) Delete(ctx context.Context, userID, todoID, version int64) *apperrors.AppError {
	deleted, err := s.repo.DeleteVersioned(ctx, todoID, userID, version)
	if err != nil {
		if err == errVersionConflict {
			existing, _ := s.repo.FindByIDAndUser(ctx, todoID, userID)
			if existing != nil {
				return apperrors.NewVersionConflictError(existing.Version)
			}
			return apperrors.NewVersionConflictError(0)
		}
		log.Printf("internal error: failed to delete todo %d for user %d: %v", todoID, userID, err)
		return apperrors.NewInternalError()
	}
	if !deleted {
		return apperrors.NewNotFoundError("todo not found")
	}
	return nil
}

// parseDate validates and parses a YYYY-MM-DD date string.
// It performs strict validation to reject invalid dates like "2024-02-30"
// which Go's time.Parse would silently normalize to "2024-03-01".
func parseDate(s string) (time.Time, error) {
	if !dateRegex.MatchString(s) {
		return time.Time{}, errInvalidDateFormat
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, errInvalidDateFormat
	}
	// Reject dates that are silently normalized (e.g. 2024-02-30 → 2024-03-01)
	if t.Format("2006-01-02") != s {
		return time.Time{}, errInvalidDateFormat
	}
	return t, nil
}

// toResponse converts a Todo model to a TodoResponse.
func toResponse(t *Todo) *TodoResponse {
	var dueDate *string
	if t.DueDate != nil {
		s := t.DueDate.Format("2006-01-02")
		dueDate = &s
	}

	return &TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     dueDate,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Version:     t.Version,
		UserID:      t.UserID,
	}
}

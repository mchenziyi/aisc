package todo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for todos.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new todo repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new todo for the given user.
func (r *Repository) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	query := `INSERT INTO todos (user_id, title, description, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`

	result := &Todo{}
	var createdAt time.Time
	var updatedAt time.Time
	var dueDate *time.Time
	err := r.pool.QueryRow(ctx, query,
		todo.UserID, todo.Title, todo.Description, todo.DueDate,
	).Scan(
		&result.ID, &result.UserID, &result.Title, &result.Description,
		&dueDate, &result.Completed, &result.Version,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}
	result.CreatedAt = formatTime(createdAt)
	result.UpdatedAt = formatTime(updatedAt)
	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		result.DueDate = &s
	}
	return result, nil
}

// GetByID retrieves a todo by ID, checking user_id for permission.
func (r *Repository) GetByID(ctx context.Context, id, userID int64) (*Todo, error) {
	query := `SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		FROM todos WHERE id = $1 AND user_id = $2`

	todo := &Todo{}
	var createdAt time.Time
	var updatedAt time.Time
	var dueDate *time.Time
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&dueDate, &todo.Completed, &todo.Version,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get todo by id: %w", err)
	}
	todo.CreatedAt = formatTime(createdAt)
	todo.UpdatedAt = formatTime(updatedAt)
	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		todo.DueDate = &s
	}
	return todo, nil
}

// GetByIDAnyUser retrieves a todo by ID without user filter (used for ownership check).
func (r *Repository) GetByIDAnyUser(ctx context.Context, id int64) (*Todo, error) {
	query := `SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		FROM todos WHERE id = $1`

	todo := &Todo{}
	var createdAt time.Time
	var updatedAt time.Time
	var dueDate *time.Time
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&dueDate, &todo.Completed, &todo.Version,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get todo by id: %w", err)
	}
	todo.CreatedAt = formatTime(createdAt)
	todo.UpdatedAt = formatTime(updatedAt)
	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		todo.DueDate = &s
	}
	return todo, nil
}

// List retrieves a paginated list of todos for a user, ordered by created_at DESC.
func (r *Repository) List(ctx context.Context, userID int64, params ListParams) (*TodoListResponse, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM todos WHERE user_id = $1`
	var total int
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count todos: %w", err)
	}

	// Calculate total pages
	totalPages := 0
	if total > 0 {
		totalPages = (total + params.PageSize - 1) / params.PageSize
	}

	// Fetch items
	offset := (params.Page - 1) * params.PageSize
	itemsQuery := `SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		FROM todos WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, itemsQuery, userID, params.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	defer rows.Close()

	items := make([]Todo, 0)
	for rows.Next() {
		var t Todo
		var createdAt time.Time
		var updatedAt time.Time
		var dueDate *time.Time
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&dueDate, &t.Completed, &t.Version,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan todo row: %w", err)
		}
		t.CreatedAt = formatTime(createdAt)
		t.UpdatedAt = formatTime(updatedAt)
		if dueDate != nil {
			s := dueDate.Format("2006-01-02")
			t.DueDate = &s
		}
		items = append(items, t)
	}

	return &TodoListResponse{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateFields holds the fields to update with explicit nullability.
type UpdateFields struct {
	Title       *string // nil = don't update
	Description  *FieldString // nil = don't update
	DueDate     *FieldString // nil = don't update
	Completed   *bool   // nil = don't update
}

// FieldString represents a nullable string field for update operations.
// Value is the string value; SetNull indicates whether to set the field to NULL.
type FieldString struct {
	Value   string
	SetNull bool
}

// UpdateWithFields performs an update with explicit field selection.
func (r *Repository) UpdateWithFields(ctx context.Context, id, userID, expectedVersion int64, fields UpdateFields) (*Todo, error) {
	query := `UPDATE todos SET version = version + 1, updated_at = NOW()`
	args := []interface{}{}
	argIdx := 1

	if fields.Title != nil {
		query += fmt.Sprintf(", title = $%d", argIdx)
		args = append(args, *fields.Title)
		argIdx++
	}

	if fields.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIdx)
		if fields.Description.SetNull {
			args = append(args, nil)
		} else {
			args = append(args, fields.Description.Value)
		}
		argIdx++
	}

	if fields.DueDate != nil {
		query += fmt.Sprintf(", due_date = $%d", argIdx)
		if fields.DueDate.SetNull {
			args = append(args, nil)
		} else {
			args = append(args, fields.DueDate.Value)
		}
		argIdx++
	}

	if fields.Completed != nil {
		query += fmt.Sprintf(", completed = $%d", argIdx)
		args = append(args, *fields.Completed)
		argIdx++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d AND version = $%d",
		argIdx, argIdx+1, argIdx+2)
	args = append(args, id, userID, expectedVersion)

	query += ` RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`

	result := &Todo{}
	var createdAt time.Time
	var updatedAt time.Time
	var dueDate *time.Time
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&result.ID, &result.UserID, &result.Title, &result.Description,
		&dueDate, &result.Completed, &result.Version,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVersionConflict
		}
		return nil, fmt.Errorf("update todo with fields: %w", err)
	}
	result.CreatedAt = formatTime(createdAt)
	result.UpdatedAt = formatTime(updatedAt)
	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		result.DueDate = &s
	}
	return result, nil
}

// UpdateCompletedOnly updates only the completed field (used for idempotent completion).
func (r *Repository) UpdateCompletedOnly(ctx context.Context, id, userID, expectedVersion int64, completed bool) (*Todo, error) {
	query := `UPDATE todos SET completed = $3, version = version + 1, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND version = $4
		RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`

	result := &Todo{}
	var createdAt time.Time
	var updatedAt time.Time
	var dueDate *time.Time
	err := r.pool.QueryRow(ctx, query, id, userID, completed, expectedVersion).Scan(
		&result.ID, &result.UserID, &result.Title, &result.Description,
		&dueDate, &result.Completed, &result.Version,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVersionConflict
		}
		return nil, fmt.Errorf("update completed only: %w", err)
	}
	result.CreatedAt = formatTime(createdAt)
	result.UpdatedAt = formatTime(updatedAt)
	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		result.DueDate = &s
	}
	return result, nil
}

// Delete removes a todo using optimistic locking.
func (r *Repository) Delete(ctx context.Context, id, userID, version int64) error {
	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2 AND version = $3`
	tag, err := r.pool.Exec(ctx, query, id, userID, version)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrVersionConflict
	}
	return nil
}

// formatTime formats a time.Time as an RFC3339 string.
func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

package todo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Todo represents a todo record from the database.
type Todo struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed"`
	Version     int64      `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UpdateFields represents the fields that can be updated on a todo.
// For each pointer field, nil means the field is not updated.
// UpdateDescription/UpdateDueDate flags indicate whether to update the field.
// If the flag is true and the Val pointer is nil, the column is set to NULL.
type UpdateFields struct {
	Title             *string
	DescriptionVal    *string    // value to set (nil means set to NULL)
	UpdateDescription bool       // whether to update the description field
	DueDateVal        *time.Time // value to set (nil means set to NULL)
	UpdateDueDate     bool       // whether to update the due_date field
	Completed         *bool
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts a new todo record.
func (r *Repository) Create(ctx context.Context, todo *Todo) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`INSERT INTO todos (user_id, title, description, due_date)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`,
		todo.UserID, todo.Title, todo.Description, todo.DueDate,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindByIDAndUser finds a todo by id and user_id (ensures data isolation).
func (r *Repository) FindByIDAndUser(ctx context.Context, id, userID int64) (*Todo, error) {
	var t Todo
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		 FROM todos WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// ListByUser returns paginated todos for a user, ordered by created_at DESC, id DESC.
func (r *Repository) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]*Todo, int, error) {
	// Count total
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM todos WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch page
	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, description, due_date, completed, version, created_at, updated_at
		 FROM todos WHERE user_id = $1
		 ORDER BY created_at DESC, id DESC
		 LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		todos = append(todos, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

// UpdateVersioned updates a todo with dynamic fields and optimistic locking.
// Returns the updated todo, nil if not found, or errVersionConflict if version mismatch.
func (r *Repository) UpdateVersioned(ctx context.Context, todoID, userID, version int64, fields *UpdateFields) (*Todo, error) {
	// Build dynamic SET clause
	var setClauses []string
	args := []interface{}{}
	argIdx := 1

	if fields.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *fields.Title)
		argIdx++
	}
	if fields.UpdateDescription {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, fields.DescriptionVal) // nil maps to SQL NULL
		argIdx++
	}
	if fields.UpdateDueDate {
		setClauses = append(setClauses, fmt.Sprintf("due_date = $%d", argIdx))
		args = append(args, fields.DueDateVal) // nil maps to SQL NULL
		argIdx++
	}
	if fields.Completed != nil {
		setClauses = append(setClauses, fmt.Sprintf("completed = $%d", argIdx))
		args = append(args, *fields.Completed)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	// Add WHERE conditions
	setClauses = append(setClauses, "version = version + 1", "updated_at = NOW()")

	sql := fmt.Sprintf(
		`UPDATE todos SET %s WHERE id = $%d AND user_id = $%d AND version = $%d
		 RETURNING id, user_id, title, description, due_date, completed, version, created_at, updated_at`,
		strings.Join(setClauses, ", "),
		argIdx, argIdx+1, argIdx+2,
	)
	args = append(args, todoID, userID, version)

	var t Todo
	err := r.pool.QueryRow(ctx, sql, args...).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description, &t.DueDate, &t.Completed, &t.Version, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Check if record exists at all
			existing, checkErr := r.FindByIDAndUser(ctx, todoID, userID)
			if checkErr != nil {
				return nil, checkErr
			}
			if existing == nil {
				return nil, nil // not found
			}
			// Record exists but version mismatch
			return nil, errVersionConflict
		}
		return nil, err
	}
	return &t, nil
}

// DeleteVersioned deletes a todo with optimistic locking.
// Returns true if deleted, false if not found, errVersionConflict if version mismatch.
func (r *Repository) DeleteVersioned(ctx context.Context, id, userID, version int64) (bool, error) {
	ct, err := r.pool.Exec(ctx,
		`DELETE FROM todos WHERE id = $1 AND user_id = $2 AND version = $3`,
		id, userID, version,
	)
	if err != nil {
		return false, err
	}

	if ct.RowsAffected() == 1 {
		return true, nil
	}

	// Check if record exists
	existing, checkErr := r.FindByIDAndUser(ctx, id, userID)
	if checkErr != nil {
		return false, checkErr
	}
	if existing == nil {
		return false, nil // not found
	}

	// Record exists but version mismatch
	return false, errVersionConflict
}

// errVersionConflict is a sentinel error for version conflicts.
var errVersionConflict = errors.New("version conflict")

package task

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, t *Task) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()

	const query = `INSERT INTO tasks(id, title, description, status, created_at) VALUES($1,$2,$3,$4,$5)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.Title, t.Description, t.Status, t.CreatedAt)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	const query = `SELECT * FROM tasks WHERE id=$1`
	t := &Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTaskNotFound
	}
	return t, err
}

func (r *Repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM tasks WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

package user_task

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Assign(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	const query = `INSERT INTO user_tasks(user_id, task_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, userID, taskID)
	return err
}

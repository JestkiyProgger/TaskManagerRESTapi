package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	Assignees   []uuid.UUID
}

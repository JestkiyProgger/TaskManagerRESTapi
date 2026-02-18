package task

import (
	"ProjectManagementAPI/internal/domain/task"
	task2 "ProjectManagementAPI/internal/repository/postgres/task"
	"context"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *task.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo RepositoryInterface
}

func NewTaskService(repo *task2.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title, description, status string, assignees []uuid.UUID) (uuid.UUID, error) {
	if title == "" {
		return uuid.Nil, task.ErrInvalidTitle
	}
	if len(assignees) == 0 {
		return uuid.Nil, task.ErrNoAssignees
	}

	t := &task.Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      status,
		Assignees:   assignees,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return uuid.Nil, err
	}

	return t.ID, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, id)
}

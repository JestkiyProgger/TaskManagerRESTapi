package task

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title, description, status string) (uuid.UUID, error) {
	if title == "" {
		return uuid.Nil, ErrInvalidTitle
	}

	t := &Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      status,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return uuid.Nil, err
	}

	return t.ID, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, id)
}

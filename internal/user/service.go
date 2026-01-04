package user

import (
	"context"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, email, name string) (uuid.UUID, error) {
	u := &User{
		ID:    uuid.New(),
		Email: email,
		Name:  name,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return uuid.Nil, err
	}

	return u.ID, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, id)
}

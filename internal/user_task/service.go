package user_task

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

func (s *Service) AssignUsersToTask(
	ctx context.Context,
	taskID uuid.UUID,
	userIDs []uuid.UUID,
) error {

	for _, userID := range userIDs {
		if err := s.repo.Assign(ctx, userID, taskID); err != nil {
			return err
		}
	}

	return nil
}

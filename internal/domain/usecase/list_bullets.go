package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

func (s *bulletService) ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
	return s.repo.List(ctx, userID)
}

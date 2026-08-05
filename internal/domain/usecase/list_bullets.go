package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
)

func (s *bulletService) ListBullets(ctx context.Context) ([]*entity.Bullet, error) {
	return s.repo.List(ctx)
}

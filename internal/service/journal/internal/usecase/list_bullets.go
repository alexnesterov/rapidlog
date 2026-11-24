package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

func (uc *Usecase) ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
	bullets, err := uc.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return bullets, nil
}

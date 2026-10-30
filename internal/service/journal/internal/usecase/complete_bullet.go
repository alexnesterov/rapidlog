package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

func (uc *usecase) CompleteBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	bullet, err := uc.repo.Get(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	err = bullet.Complete()
	if errors.Is(err, entity.ErrBulletAlreadyCompleted) {
		return bullet, nil
	}
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

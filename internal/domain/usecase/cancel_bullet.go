package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

func (s *bulletService) CancelBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	bullet, err := s.repo.Get(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	err = bullet.Cancel()
	if errors.Is(err, entity.ErrBulletAlreadyCancelled) {
		return bullet, nil
	}
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

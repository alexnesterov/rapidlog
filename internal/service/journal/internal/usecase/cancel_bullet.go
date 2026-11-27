package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
)

type CancelBullet struct {
	repo BulletRepository
}

func NewCancelBullet(repo BulletRepository) *CancelBullet {
	return &CancelBullet{repo: repo}
}

type CancelBulletInput struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (uc *CancelBullet) Execute(ctx context.Context, in CancelBulletInput) (*domain.Bullet, error) {
	bullet, err := uc.repo.Get(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}

	err = bullet.Cancel()
	if errors.Is(err, domain.ErrBulletAlreadyCancelled) {
		return bullet, nil
	}
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

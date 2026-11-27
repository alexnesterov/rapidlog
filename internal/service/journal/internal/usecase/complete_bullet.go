package usecase

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
)

type CompleteBullet struct {
	repo BulletRepository
}

func NewCompleteBullet(repo BulletRepository) *CompleteBullet {
	return &CompleteBullet{repo: repo}
}

type CompleteBulletInput struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (uc *CompleteBullet) Execute(ctx context.Context, in CompleteBulletInput) (*domain.Bullet, error) {
	bullet, err := uc.repo.Get(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}

	err = bullet.Complete()
	if errors.Is(err, domain.ErrBulletAlreadyCompleted) {
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

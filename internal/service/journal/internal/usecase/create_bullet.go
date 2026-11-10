package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

type CreateBulletInput struct {
	UserID  uuid.UUID
	Kind    entity.Kind
	Content string
}

func (uc *Usecase) CreateBullet(ctx context.Context, input CreateBulletInput) (*entity.Bullet, error) {
	bullet, err := entity.NewBullet(input.UserID, input.Kind, input.Content)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

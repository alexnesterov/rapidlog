package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
)

type CreateBullet struct {
	repo BulletRepository
}

func NewCreateBullet(repo BulletRepository) *CreateBullet {
	return &CreateBullet{repo: repo}
}

type CreateBulletInput struct {
	UserID  uuid.UUID
	Kind    domain.Kind
	Content string
}

func (uc *CreateBullet) Execute(ctx context.Context, in CreateBulletInput) (*domain.Bullet, error) {
	bullet, err := domain.NewBullet(in.UserID, in.Kind, in.Content)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

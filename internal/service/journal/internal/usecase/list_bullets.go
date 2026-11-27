package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/google/uuid"
)

type ListBullets struct {
	repo BulletRepository
}

func NewListBullets(repo BulletRepository) *ListBullets {
	return &ListBullets{repo: repo}
}

type ListBulletsInput struct {
	UserID uuid.UUID
}

func (uc *ListBullets) Execute(ctx context.Context, in ListBulletsInput) ([]*domain.Bullet, error) {
	bullets, err := uc.repo.List(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return bullets, nil
}

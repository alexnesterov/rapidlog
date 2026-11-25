package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

type MigrateBullet struct {
	repo BulletRepository
}

func NewMigrateBullet(repo BulletRepository) *MigrateBullet {
	return &MigrateBullet{repo: repo}
}

type MigrateBulletInput struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (uc *MigrateBullet) Execute(ctx context.Context, in MigrateBulletInput) (*entity.Bullet, error) {
	bullet, err := uc.repo.Get(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}

	err = bullet.Migrate()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

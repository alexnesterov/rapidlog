package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

func (uc *usecase) MigrateBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	bullet, err := uc.repo.Get(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	err = bullet.Migrate()
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

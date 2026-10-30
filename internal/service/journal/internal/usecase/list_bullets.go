package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

func (uc *usecase) ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
	bullets, err := uc.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(bullets) > 0 {
		return bullets, nil
	}

	bullets, err = entity.NewDemoBullets(userID)
	if err != nil {
		return nil, err
	}

	err = uc.transactor.WithTx(ctx, func(ctx context.Context) error {
		for _, b := range bullets {
			err := uc.repo.Create(ctx, b)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return bullets, nil
}

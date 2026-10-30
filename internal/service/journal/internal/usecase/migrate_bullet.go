package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

func (uc *usecase) MigrateBullet(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error) {
	var migrated *entity.Bullet

	err := uc.transactor.WithTx(ctx, func(txCtx context.Context) error {
		bullet, err := uc.repo.Get(txCtx, id, userID)
		if err != nil {
			return err
		}

		migrated, err = bullet.Migrate()
		if err != nil {
			return err
		}

		if err := uc.repo.Update(txCtx, bullet); err != nil {
			return err
		}

		if err := uc.repo.Create(txCtx, migrated); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return migrated, nil
}

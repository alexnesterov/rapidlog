package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

func (s *bulletService) MigrateBullet(ctx context.Context, id uuid.UUID) (*entity.Bullet, error) {
	var migrated *entity.Bullet

	err := s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
		bullet, err := s.repo.Get(txCtx, id)
		if err != nil {
			return err
		}

		migrated, err = bullet.Migrate()
		if err != nil {
			return err
		}

		if err := s.repo.Update(txCtx, bullet); err != nil {
			return err
		}

		if err := s.repo.Create(txCtx, migrated); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return migrated, nil
}

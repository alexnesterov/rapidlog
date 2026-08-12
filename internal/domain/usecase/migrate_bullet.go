package usecase

import (
	"context"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

func (s *bulletService) MigrateBullet(ctx context.Context, id uuid.UUID) (*entity.Bullet, error) {
	bullet, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := bullet.ValidateMigrate(); err != nil {
		return nil, err
	}

	now := time.Now()
	bullet.Signifier = entity.SignifierMigrated
	bullet.UpdatedAt = now

	if err := s.repo.Update(ctx, bullet); err != nil {
		return nil, err
	}

	migrated := &entity.Bullet{
		ID:        uuid.New(),
		Type:      bullet.Type,
		Signifier: entity.SignifierOpen,
		Content:   bullet.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, migrated); err != nil {
		return nil, err
	}

	return migrated, nil
}

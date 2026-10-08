package usecase

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
)

func (s *bulletService) CreateBullet(ctx context.Context, input port.CreateBulletInput) (*entity.Bullet, error) {
	bullet, err := entity.NewBullet(input.UserID, input.Type, input.Content)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

package usecase

import (
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

func (s *bulletService) CreateBullet(input port.CreateBulletInput) (*entity.Bullet, error) {
	now := time.Now()

	bulletType := input.Type
	if bulletType == "" {
		bulletType = entity.BulletTask
	}

	bullet := &entity.Bullet{
		ID:        uuid.New(),
		Type:      bulletType,
		Signifier: entity.SignifierOpen,
		Content:   input.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := bullet.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

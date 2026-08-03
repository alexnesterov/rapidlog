package usecase

import (
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

func (s *bulletService) CompleteBullet(id uuid.UUID) (*entity.Bullet, error) {
	bullet, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if bullet.Signifier == entity.SignifierCompleted {
		return bullet, nil
	}

	bullet.Signifier = entity.SignifierCompleted
	bullet.UpdatedAt = time.Now()

	if err := s.repo.Update(bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

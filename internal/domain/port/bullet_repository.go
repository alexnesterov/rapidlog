package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
)

type BulletRepository interface {
	Create(bullet *entity.Bullet) error
	List() ([]*entity.Bullet, error)
	Get(id uuid.UUID) (*entity.Bullet, error)
	Update(bullet *entity.Bullet) error
}

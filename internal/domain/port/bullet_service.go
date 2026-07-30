package port

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
)

type CreateBulletRequest struct {
	Title string
}

type BulletService interface {
	CreateBullet(req CreateBulletRequest) (*entity.Bullet, error)
	ListBullets() ([]*entity.Bullet, error)
}

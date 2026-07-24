package usecase

import (
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

type BulletService struct {
	repo port.BulletRepository
}

func NewBulletService(r port.BulletRepository) *BulletService {
	return &BulletService{
		repo: r,
	}
}

func (s *BulletService) CreateBullet(req port.CreateBulletRequest) (*entity.Bullet, error) {
	now := time.Now()

	bullet := &entity.Bullet{
		ID:        uuid.New(),
		Title:     req.Title,
		Status:    entity.StatusOpen,
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

func (s *BulletService) ListBullets() ([]*entity.Bullet, error) {
	return nil, nil
}

func (s *BulletService) ReadBullet(id uuid.UUID) (*entity.Bullet, error) {
	return nil, nil
}

func (s *BulletService) UpdateBullet(req port.UpdateBulletRequest) (*entity.Bullet, error) {
	return nil, nil
}

func (s *BulletService) DeleteBullet(id uuid.UUID) error {
	return nil
}

var _ port.BulletService = &BulletService{}

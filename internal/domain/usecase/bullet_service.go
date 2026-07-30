package usecase

import (
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

type bulletService struct {
	repo port.BulletRepository
}

func NewBulletService(r port.BulletRepository) *bulletService {
	return &bulletService{
		repo: r,
	}
}

func (s *bulletService) CreateBullet(req port.CreateBulletRequest) (*entity.Bullet, error) {
	now := time.Now()

	bullet := &entity.Bullet{
		ID:        uuid.New(),
		Title:     req.Title,
		Status:    entity.StatusOpened,
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

func (s *bulletService) ListBullets() ([]*entity.Bullet, error) {
	return s.repo.List()
}

func (s *bulletService) CompleteBullet(id uuid.UUID) (*entity.Bullet, error) {
	bullet, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if bullet.Status == entity.StatusCompleted {
		return bullet, nil
	}

	bullet.Status = entity.StatusCompleted
	bullet.UpdatedAt = time.Now()

	if err := s.repo.Update(bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}

var _ port.BulletService = &bulletService{}

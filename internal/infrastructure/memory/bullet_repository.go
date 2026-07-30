package memory

import (
	"sync"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

type bulletRepository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*entity.Bullet
}

func NewBulletRepository() *bulletRepository {
	return &bulletRepository{
		data: make(map[uuid.UUID]*entity.Bullet),
	}
}

func (r *bulletRepository) Create(bullet *entity.Bullet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	bulletCopy := *bullet
	r.data[bullet.ID] = &bulletCopy

	return nil
}

func (r *bulletRepository) List() ([]*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullets := make([]*entity.Bullet, 0, len(r.data))
	for _, bullet := range r.data {
		bulletCopy := *bullet
		bullets = append(bullets, &bulletCopy)
	}

	return bullets, nil
}

func (r *bulletRepository) Get(id uuid.UUID) (*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullet, ok := r.data[id]
	if !ok {
		return nil, port.ErrNotFound
	}

	bulletCopy := *bullet
	return &bulletCopy, nil
}

var _ port.BulletRepository = &bulletRepository{}

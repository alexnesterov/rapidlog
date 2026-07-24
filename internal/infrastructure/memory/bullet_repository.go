package memory

import (
	"sync"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

type bulletMemory struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*entity.Bullet
}

func NewMemoryBulletRepository() *bulletMemory {
	return &bulletMemory{
		data: make(map[uuid.UUID]*entity.Bullet),
	}
}

func (r *bulletMemory) Create(bullet *entity.Bullet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	bulletCopy := *bullet
	r.data[bullet.ID] = &bulletCopy

	return nil
}

func (r *bulletMemory) Read(id uuid.UUID) (*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullet, ok := r.data[id]
	if !ok {
		return nil, port.ErrNotFound
	}

	bulletCopy := *bullet
	return &bulletCopy, nil
}

func (r *bulletMemory) List() ([]*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullets := make([]*entity.Bullet, 0, len(r.data))
	for _, bullet := range r.data {
		bulletCopy := *bullet
		bullets = append(bullets, &bulletCopy)
	}

	return bullets, nil
}

func (r *bulletMemory) Update(bullet *entity.Bullet) error {
	return nil
}

func (r *bulletMemory) Delete(id uuid.UUID) error {
	return nil
}

var _ port.BulletRepository = &bulletMemory{}

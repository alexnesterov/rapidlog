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
	r.data[bullet.ID] = bullet
	return nil
}

func (r *bulletMemory) Read(id uuid.UUID) (*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullet, ok := r.data[id]
	if !ok {
		return nil, port.ErrNotFound
	}

	return bullet, nil
}

// var _ port.BulletRepository = &bulletMemory{}

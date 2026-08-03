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

func (r *bulletRepository) store(bullet *entity.Bullet) {
	bulletCopy := *bullet
	r.data[bullet.ID] = &bulletCopy
}

func (r *bulletRepository) clone(bullet *entity.Bullet) *entity.Bullet {
	bulletCopy := *bullet
	return &bulletCopy
}

func (r *bulletRepository) find(id uuid.UUID) (*entity.Bullet, error) {
	bullet, ok := r.data[id]
	if !ok {
		return nil, port.ErrNotFound
	}

	return r.clone(bullet), nil
}

func (r *bulletRepository) Create(bullet *entity.Bullet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store(bullet)

	return nil
}

func (r *bulletRepository) List() ([]*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bullets := make([]*entity.Bullet, 0, len(r.data))
	for _, bullet := range r.data {
		bullets = append(bullets, r.clone(bullet))
	}

	return bullets, nil
}

func (r *bulletRepository) Get(id uuid.UUID) (*entity.Bullet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.find(id)
}

func (r *bulletRepository) Update(bullet *entity.Bullet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := r.find(bullet.ID); err != nil {
		return err
	}

	r.store(bullet)

	return nil
}

var _ port.BulletRepository = &bulletRepository{}

package memory

import (
	"sync"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulletMemory_CreateAndRead(t *testing.T) {
	repo := NewMemoryBulletRepository()

	bullet := &entity.Bullet{
		ID:    uuid.New(),
		Title: "Заголовок",
	}

	err := repo.Create(bullet)
	require.NoError(t, err)

	got, err := repo.Read(bullet.ID)
	require.NoError(t, err)
	assert.Equal(t, bullet, got)
}

func TestBulletMemory_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryBulletRepository()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			_ = repo.Create(&entity.Bullet{ID: uuid.New(), Title: "x"})
		})
	}

	wg.Wait()
}

func TestBulletMemory_Read_NotFound(t *testing.T) {
	repo := NewMemoryBulletRepository()

	got, err := repo.Read(uuid.New())
	assert.Nil(t, got)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

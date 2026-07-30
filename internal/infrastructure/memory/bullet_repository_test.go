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

func TestBulletMemory_List(t *testing.T) {
	b1 := &entity.Bullet{ID: uuid.New(), Title: "Первый"}
	b2 := &entity.Bullet{ID: uuid.New(), Title: "Второй"}

	cases := []struct {
		name string
		seed []*entity.Bullet
		want []*entity.Bullet
	}{
		{
			name: "with items",
			seed: []*entity.Bullet{b1, b2},
			want: []*entity.Bullet{b1, b2},
		},
		{
			name: "empty",
			seed: []*entity.Bullet{},
			want: []*entity.Bullet{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewBulletRepository()

			for _, b := range tc.seed {
				require.NoError(t, repo.Create(b))
			}

			got, err := repo.List()
			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.ElementsMatch(t, tc.want, got)
		})
	}
}

func TestBulletMemory_List_ReturnsCopies(t *testing.T) {
	repo := NewBulletRepository()

	bullet := &entity.Bullet{ID: uuid.New(), Title: "Заголовок"}
	require.NoError(t, repo.Create(bullet))

	got, err := repo.List()
	require.NoError(t, err)

	got[0].Title = "Изменено"

	again, err := repo.List()
	require.NoError(t, err)
	assert.Equal(t, "Заголовок", again[0].Title)
}

func TestBulletMemory_CreateAndRead(t *testing.T) {
	repo := NewBulletRepository()

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
	repo := NewBulletRepository()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			_ = repo.Create(&entity.Bullet{ID: uuid.New(), Title: "x"})
		})
	}

	wg.Wait()
}

func TestBulletMemory_Create_StoresCopy(t *testing.T) {
	repo := NewBulletRepository()

	bullet := &entity.Bullet{ID: uuid.New(), Title: "Заголовок"}
	require.NoError(t, repo.Create(bullet))

	bullet.Title = "Изменено"

	got, err := repo.Read(bullet.ID)
	require.NoError(t, err)
	assert.Equal(t, "Заголовок", got.Title)
}

func TestBulletMemory_Read_NotFound(t *testing.T) {
	repo := NewBulletRepository()

	got, err := repo.Read(uuid.New())
	assert.Nil(t, got)
	assert.ErrorIs(t, err, port.ErrNotFound)
}

func TestBulletMemory_Read_ReturnsCopy(t *testing.T) {
	repo := NewBulletRepository()

	bullet := &entity.Bullet{ID: uuid.New(), Title: "Заголовок"}
	require.NoError(t, repo.Create(bullet))

	got, err := repo.Read(bullet.ID)
	require.NoError(t, err)

	got.Title = "Изменено"

	again, err := repo.Read(bullet.ID)
	require.NoError(t, err)
	assert.Equal(t, "Заголовок", again.Title)
}

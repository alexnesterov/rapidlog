package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBullet(t *testing.T) {
	var errRepo = errors.New("repo error")

	cases := []struct {
		name      string
		req       port.CreateBulletRequest
		setupMock func(m *mocks.MockBulletRepository)
		want      *entity.Bullet
		wantErr   error
	}{
		{
			name: "success",
			req:  port.CreateBulletRequest{Title: "Заголовок"},
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().
					Create(mock.AnythingOfType("*entity.Bullet")).
					Return(nil).
					Once()
			},
			want:    &entity.Bullet{Title: "Заголовок"},
			wantErr: nil,
		},
		{
			name: "repo error",
			req:  port.CreateBulletRequest{Title: "Заголовок"},
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().
					Create(mock.AnythingOfType("*entity.Bullet")).
					Return(errRepo).
					Once()
			},
			want:    nil,
			wantErr: errRepo,
		},
		{
			name:      "empty title",
			req:       port.CreateBulletRequest{Title: ""},
			setupMock: func(m *mocks.MockBulletRepository) {},
			want:      nil,
			wantErr:   entity.ErrTitleRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockBulletRepo := mocks.NewMockBulletRepository(t)
			tc.setupMock(mockBulletRepo)
			uc := NewBulletService(mockBulletRepo)

			got, err := uc.CreateBullet(tc.req)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.want.Title, got.Title)
			assert.NotEqual(t, uuid.Nil, got.ID)
			assert.Equal(t, entity.SignifierOpen, got.Signifier)
			assert.False(t, got.CreatedAt.IsZero())
			assert.False(t, got.UpdatedAt.IsZero())
		})
	}
}

func TestListBullets(t *testing.T) {
	want := []*entity.Bullet{{Title: "Заголовок"}}

	mockRepo := mocks.NewMockBulletRepository(t)
	mockRepo.EXPECT().
		List().
		Return(want, nil).
		Once()

	uc := NewBulletService(mockRepo)

	got, err := uc.ListBullets()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCompleteBullet(t *testing.T) {
	oldUpdatedAt := time.Now().Add(-24 * time.Hour)

	cases := []struct {
		name                 string
		setupMock            func(m *mocks.MockBulletRepository)
		wantUpdatedAtChanged bool
		wantErr              error
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().Get(mock.AnythingOfType("uuid.UUID")).
					Return(&entity.Bullet{Signifier: entity.SignifierOpen, UpdatedAt: oldUpdatedAt}, nil).
					Once()
				m.EXPECT().Update(mock.AnythingOfType("*entity.Bullet")).
					Return(nil).
					Once()
			},
			wantUpdatedAtChanged: true,
		},
		{
			name: "already completed",
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().Get(mock.AnythingOfType("uuid.UUID")).
					Return(&entity.Bullet{Signifier: entity.SignifierCompleted, UpdatedAt: oldUpdatedAt}, nil).
					Once()
			},
		},
		{
			name: "not found",
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().Get(mock.AnythingOfType("uuid.UUID")).
					Return(nil, port.ErrNotFound).
					Once()
			},
			wantErr: port.ErrNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewMockBulletRepository(t)
			tc.setupMock(mockRepo)

			uc := NewBulletService(mockRepo)

			got, err := uc.CompleteBullet(uuid.New())
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, entity.SignifierCompleted, got.Signifier)

			if tc.wantUpdatedAtChanged {
				assert.True(t, got.UpdatedAt.After(oldUpdatedAt))
			}

			if !tc.wantUpdatedAtChanged {
				assert.Equal(t, oldUpdatedAt, got.UpdatedAt)
			}
		})
	}
}

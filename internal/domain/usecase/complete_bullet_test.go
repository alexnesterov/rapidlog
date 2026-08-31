package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCompleteBulletUseCase(t *testing.T) {
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
				m.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return(&entity.Bullet{Signifier: entity.SignifierOpen, UpdatedAt: oldUpdatedAt}, nil).
					Once()
				m.EXPECT().Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
					Return(nil).
					Once()
			},
			wantUpdatedAtChanged: true,
		},
		{
			name: "already completed",
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return(&entity.Bullet{Signifier: entity.SignifierCompleted, UpdatedAt: oldUpdatedAt}, nil).
					Once()
			},
		},
		{
			name: "not found",
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
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

			mockTxMgr := mocks.NewMockTransactionManager(t)

			uc := usecase.NewBulletService(mockRepo, mockTxMgr)

			got, err := uc.CompleteBullet(context.Background(), uuid.New(), uuid.New())
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

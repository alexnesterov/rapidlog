package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBulletUseCase(t *testing.T) {
	var errRepo = errors.New("repo error")

	cases := []struct {
		name      string
		input     port.CreateBulletInput
		setupMock func(m *mocks.MockBulletRepository)
		want      *entity.Bullet
		wantErr   error
	}{
		{
			name: "success",
			input: port.CreateBulletInput{
				Content: "Заголовок",
				Type:    entity.BulletTask,
				UserID:  uuid.New(),
			},
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
					Return(nil).
					Once()
			},
			want:    &entity.Bullet{Content: "Заголовок", Type: entity.BulletTask},
			wantErr: nil,
		},
		{
			name: "empty type",
			input: port.CreateBulletInput{
				Type:    "",
				Content: "Заголовок",
			},
			setupMock: func(m *mocks.MockBulletRepository) {},
			want:      nil,
			wantErr:   entity.ErrTypeInvalid,
		},
		{
			name: "repo error",
			input: port.CreateBulletInput{
				Type:    entity.BulletTask,
				Content: "Заголовок",
				UserID:  uuid.New(),
			},
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
					Return(errRepo).
					Once()
			},
			want:    nil,
			wantErr: errRepo,
		},
		{
			name: "empty content",
			input: port.CreateBulletInput{
				Type:    entity.BulletTask,
				Content: "",
			},
			setupMock: func(m *mocks.MockBulletRepository) {},
			want:      nil,
			wantErr:   entity.ErrContentRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewMockBulletRepository(t)
			tc.setupMock(mockRepo)
			mockTxMgr := mocks.NewMockTransactionManager(t)
			uc := usecase.NewBulletService(mockRepo, mockTxMgr)

			got, err := uc.CreateBullet(context.Background(), tc.input)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.want.Content, got.Content)
			assert.NotEqual(t, uuid.Nil, got.ID)
			assert.Equal(t, tc.want.Type, got.Type)
			assert.Equal(t, entity.SignifierOpen, got.Signifier)
			assert.False(t, got.CreatedAt.IsZero())
			assert.False(t, got.UpdatedAt.IsZero())
		})
	}
}

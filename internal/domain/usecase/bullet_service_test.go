package usecase

import (
	"errors"
	"testing"

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
			assert.Equal(t, entity.StatusOpen, got.Status)
			assert.False(t, got.CreatedAt.IsZero())
			assert.False(t, got.UpdatedAt.IsZero())
		})
	}
}

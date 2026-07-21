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
		name    string
		req     port.CreateBulletRequest
		repoErr error
		want    *entity.Bullet
		wantErr error
	}{
		{
			name:    "success",
			req:     port.CreateBulletRequest{Title: "Заголовок"},
			repoErr: nil,
			want:    &entity.Bullet{Title: "Заголовок"},
			wantErr: nil,
		},
		{
			name:    "repo error",
			req:     port.CreateBulletRequest{Title: "Заголовок"},
			repoErr: errRepo,
			want:    nil,
			wantErr: errRepo,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewMockBulletRepository(t)

			mockRepo.EXPECT().
				Create(mock.AnythingOfType("*entity.Bullet")).
				Return(tc.repoErr).
				Once()

			uc := &BulletService{
				BulletRepo: mockRepo,
			}

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
		})
	}
}

package usecase_test

import (
	"errors"
	"testing"

	mocksTransactor "github.com/alexnesterov/rapidlog-api/internal/pkg/transactor/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
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
			name: "empty type",
			input: port.CreateBulletInput{
				Kind:    "",
				Content: "Заголовок",
			},
			setupMock: func(m *mocks.MockBulletRepository) {},
			want:      nil,
			wantErr:   entity.ErrBulletKindInvalid,
		},
		{
			name: "repo error",
			input: port.CreateBulletInput{
				Kind:    entity.KindTask,
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
				Kind:    entity.KindTask,
				Content: "",
			},
			setupMock: func(m *mocks.MockBulletRepository) {},
			want:      nil,
			wantErr:   entity.ErrBulletContentRequired,
		},
		{
			name: "success",
			input: port.CreateBulletInput{
				Content: "Заголовок",
				Kind:    entity.KindTask,
				UserID:  uuid.New(),
			},
			setupMock: func(m *mocks.MockBulletRepository) {
				m.EXPECT().
					Create(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
					Return(nil).
					Once()
			},
			want:    &entity.Bullet{Content: "Заголовок", Kind: entity.KindTask},
			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewMockBulletRepository(t)
			tc.setupMock(mockRepo)
			mockTxMgr := mocksTransactor.NewMockTransactor(t)
			uc := usecase.New(mockRepo, mockTxMgr)

			got, gotErr := uc.CreateBullet(t.Context(), tc.input)
			if tc.wantErr != nil {
				require.Error(t, gotErr)
				assert.Nil(t, got)
				assert.ErrorIs(t, gotErr, tc.wantErr)

				return
			}

			require.NoError(t, gotErr)

			assert.Equal(t, tc.want.Content, got.Content)
			assert.NotEqual(t, uuid.Nil, got.ID)
			assert.Equal(t, tc.want.Kind, got.Kind)
			assert.Equal(t, entity.SignifierOpen, got.Signifier)
			assert.False(t, got.CreatedAt.IsZero())
			assert.False(t, got.UpdatedAt.IsZero())
		})
	}
}

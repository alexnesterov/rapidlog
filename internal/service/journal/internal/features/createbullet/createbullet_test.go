package createbullet_test

import (
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/createbullet"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/createbullet/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBulletUseCase(t *testing.T) {
	var errRepo = errors.New("repo error")

	cases := []struct {
		name      string
		input     createbullet.Input
		setupMock func(m *mocks.MockStore)
		want      *entity.Bullet
		wantErr   error
	}{
		{
			name: "empty type",
			input: createbullet.Input{
				Kind:    "",
				Content: "Заголовок",
			},
			setupMock: func(m *mocks.MockStore) {},
			want:      nil,
			wantErr:   entity.ErrBulletKindInvalid,
		},
		{
			name: "repo error",
			input: createbullet.Input{
				Kind:    entity.KindTask,
				Content: "Заголовок",
				UserID:  uuid.New(),
			},
			setupMock: func(m *mocks.MockStore) {
				m.EXPECT().
					Insert(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
					Return(errRepo).
					Once()
			},
			want:    nil,
			wantErr: errRepo,
		},
		{
			name: "empty content",
			input: createbullet.Input{
				Kind:    entity.KindTask,
				Content: "",
			},
			setupMock: func(m *mocks.MockStore) {},
			want:      nil,
			wantErr:   entity.ErrBulletContentRequired,
		},
		{
			name: "success",
			input: createbullet.Input{
				Content: "Заголовок",
				Kind:    entity.KindTask,
				UserID:  uuid.New(),
			},
			setupMock: func(m *mocks.MockStore) {
				m.EXPECT().
					Insert(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
					Return(nil).
					Once()
			},
			want:    &entity.Bullet{Content: "Заголовок", Kind: entity.KindTask},
			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewMockStore(t)
			tc.setupMock(mockRepo)
			uc := createbullet.New(mockRepo)

			got, gotErr := uc.Execute(t.Context(), tc.input)
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

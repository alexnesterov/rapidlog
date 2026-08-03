package usecase_test

import (
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBulletsUseCase(t *testing.T) {
	want := []*entity.Bullet{{Content: "Заголовок"}}

	mockRepo := mocks.NewMockBulletRepository(t)
	mockRepo.EXPECT().
		List().
		Return(want, nil).
		Once()

	uc := usecase.NewBulletService(mockRepo)

	got, err := uc.ListBullets()
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

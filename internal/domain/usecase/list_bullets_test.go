package usecase_test

import (
	"context"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListBulletsUseCase(t *testing.T) {
	userID := uuid.New()
	want := []*entity.Bullet{{Content: "Заголовок", UserID: userID}}

	mockRepo := mocks.NewMockBulletRepository(t)
	mockRepo.EXPECT().
		List(mock.Anything, userID).
		Return(want, nil).
		Once()

	mockTxMgr := mocks.NewMockTransactionManager(t)

	uc := usecase.NewBulletService(mockRepo, mockTxMgr)

	got, err := uc.ListBullets(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

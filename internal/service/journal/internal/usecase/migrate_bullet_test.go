package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mocksTransactor "github.com/alexnesterov/rapidlog-api/internal/pkg/transactor/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MigrateBulletUseCaseSuite struct {
	suite.Suite
}

func TestMigrateBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(MigrateBulletUseCaseSuite))
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotFound() {
	id := uuid.New()
	userID := uuid.New()

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id, userID).
		Return(nil, port.ErrNotFound).
		Once()

	mockTxMgr := mocksTransactor.NewMockTransactor(s.T())

	uc := usecase.New(mockRepo, mockTxMgr)

	_, err := uc.MigrateBullet(context.Background(), id, userID)
	s.Error(err)
	s.ErrorIs(err, port.ErrNotFound)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_UpdateError() {
	id := uuid.New()
	userID := uuid.New()
	wantErr := errors.New("update failed")

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierOpen,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	mockTxMgr := mocksTransactor.NewMockTransactor(s.T())

	uc := usecase.New(mockRepo, mockTxMgr)

	got, err := uc.MigrateBullet(context.Background(), id, userID)
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotOpenTask() {
	id := uuid.New()
	userID := uuid.New()

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierCompleted,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	mockTxMgr := mocksTransactor.NewMockTransactor(s.T())

	uc := usecase.New(mockRepo, mockTxMgr)

	got, err := uc.MigrateBullet(context.Background(), id, userID)
	s.Nil(got)
	s.ErrorIs(err, entity.ErrBulletNotOpenTask)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_Success() {
	id := uuid.New()
	userID := uuid.New()
	dayAgo := time.Now().Add(-24 * time.Hour)

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierOpen,
		Content:   "Test task",
		CreatedAt: dayAgo,
		UpdatedAt: dayAgo,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	mockRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(b *entity.Bullet) bool {
		return b.ID == bullet.ID &&
			b.CreatedAt.After(dayAgo) &&
			b.UpdatedAt.After(dayAgo)
	})).
		Return(nil).
		Once()

	mockTxMgr := mocksTransactor.NewMockTransactor(s.T())

	uc := usecase.New(mockRepo, mockTxMgr)

	got, err := uc.MigrateBullet(s.T().Context(), id, userID)

	s.NoError(err)
	s.Same(bullet, got)
	s.True(got.CreatedAt.After(dayAgo))
	s.True(got.UpdatedAt.After(dayAgo))
}

package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
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

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_Success() {
	id := uuid.New()
	initialUpdatedAt := time.Now().Add(-24 * time.Hour)

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierOpen,
		Content:   "Test task",
		UpdatedAt: initialUpdatedAt,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id).
		Return(bullet, nil).
		Once()

	mockRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(b *entity.Bullet) bool {
		return b.ID == bullet.ID &&
			b.Signifier == entity.SignifierMigrated &&
			b.UpdatedAt.After(initialUpdatedAt)
	})).
		Return(nil).
		Once()

	mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(b *entity.Bullet) bool {
		return b.ID != bullet.ID &&
			b.Type == entity.BulletTask &&
			b.Signifier == entity.SignifierOpen &&
			b.Content == bullet.Content
	})).
		Return(nil).
		Once()

	uc := usecase.NewBulletService(mockRepo)

	got, err := uc.MigrateBullet(context.Background(), id)

	s.NoError(err)
	s.NotEqual(bullet.ID, got.ID)
	s.Equal(bullet.Content, got.Content)
	s.Equal(entity.BulletTask, got.Type)
	s.Equal(entity.SignifierOpen, got.Signifier)
	s.True(got.UpdatedAt.After(initialUpdatedAt))
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotFound() {
	id := uuid.New()

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id).
		Return(nil, port.ErrNotFound).
		Once()

	uc := usecase.NewBulletService(mockRepo)

	_, err := uc.MigrateBullet(context.Background(), id)
	s.Error(err)
	s.ErrorIs(err, port.ErrNotFound)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_UpdateError() {
	id := uuid.New()
	wantErr := errors.New("update failed")

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierOpen,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id).
		Return(bullet, nil).
		Once()

	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	uc := usecase.NewBulletService(mockRepo)

	got, err := uc.MigrateBullet(context.Background(), id)
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_CreateError() {
	id := uuid.New()
	wantErr := errors.New("create failed")

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierOpen,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id).
		Return(bullet, nil).
		Once()

	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	uc := usecase.NewBulletService(mockRepo)

	got, err := uc.MigrateBullet(context.Background(), id)
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotOpenTask() {
	id := uuid.New()

	bullet := &entity.Bullet{
		ID:        id,
		Type:      entity.BulletTask,
		Signifier: entity.SignifierCompleted,
	}

	mockRepo := mocks.NewMockBulletRepository(s.T())
	mockRepo.EXPECT().Get(mock.Anything, id).
		Return(bullet, nil).
		Once()

	uc := usecase.NewBulletService(mockRepo)

	got, err := uc.MigrateBullet(context.Background(), id)
	s.Nil(got)
	s.ErrorIs(err, entity.ErrNotOpenTask)
}

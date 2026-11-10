package usecase_test

import (
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
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocksTransactor.MockTransactor
	uc             port.Usecase
}

func TestMigrateBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(MigrateBulletUseCaseSuite))
}

func (s *MigrateBulletUseCaseSuite) SetupTest() {
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocksTransactor.NewMockTransactor(s.T())
	s.uc = usecase.New(s.mockBulletRepo, s.mockTxMgr)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotFound() {
	id := uuid.New()
	userID := uuid.New()

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, id, userID).
		Return(nil, entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.MigrateBullet(s.T().Context(), id, userID)
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrNotFound)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_UpdateError() {
	id := uuid.New()
	userID := uuid.New()
	wantErr := errors.New("update failed")
	dayAgo := time.Now().Add(-24 * time.Hour)

	bullet := &entity.Bullet{
		ID:        id,
		Kind:      entity.KindTask,
		Signifier: entity.SignifierOpen,
		CreatedAt: dayAgo,
		UpdatedAt: dayAgo,
	}

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.MigrateBullet(s.T().Context(), id, userID)
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotOpenTask() {
	id := uuid.New()
	userID := uuid.New()

	bullet := &entity.Bullet{
		ID:        id,
		Kind:      entity.KindTask,
		Signifier: entity.SignifierCompleted,
	}

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	got, gotErr := s.uc.MigrateBullet(s.T().Context(), id, userID)
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrBulletNotOpenTask)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_Success() {
	id := uuid.New()
	userID := uuid.New()
	dayAgo := time.Now().Add(-24 * time.Hour)

	bullet := &entity.Bullet{
		ID:        id,
		Kind:      entity.KindTask,
		Signifier: entity.SignifierOpen,
		Content:   "Test task",
		CreatedAt: dayAgo,
		UpdatedAt: dayAgo,
	}

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.MigrateBullet(s.T().Context(), id, userID)

	s.NoError(gotErr)
	s.Same(bullet, got)
	s.True(got.CreatedAt.After(dayAgo))
	s.True(got.UpdatedAt.After(dayAgo))
}

package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MigrateBulletUseCaseSuite struct {
	suite.Suite
	bulletRepoMock *mocks.MockBulletRepository
	uc             *usecase.MigrateBullet
}

func TestMigrateBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(MigrateBulletUseCaseSuite))
}

func (s *MigrateBulletUseCaseSuite) SetupTest() {
	s.bulletRepoMock = mocks.NewMockBulletRepository(s.T())
	s.uc = usecase.NewMigrateBullet(s.bulletRepoMock)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotFound() {
	id := uuid.New()
	userID := uuid.New()

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, id, userID).
		Return(nil, domain.ErrNotFound).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.MigrateBulletInput{
		ID:     id,
		UserID: userID,
	})
	s.Nil(got)
	s.ErrorIs(gotErr, domain.ErrNotFound)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_SaveError() {
	id := uuid.New()
	userID := uuid.New()
	wantErr := errors.New("update failed")
	dayAgo := time.Now().Add(-24 * time.Hour)

	bullet := &domain.Bullet{
		ID:        id,
		Kind:      domain.KindTask,
		Signifier: domain.SignifierOpen,
		CreatedAt: dayAgo,
		UpdatedAt: dayAgo,
	}

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*domain.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.MigrateBulletInput{
		ID:     id,
		UserID: userID,
	})
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_NotOpenTask() {
	id := uuid.New()
	userID := uuid.New()

	bullet := &domain.Bullet{
		ID:        id,
		Kind:      domain.KindTask,
		Signifier: domain.SignifierCompleted,
	}

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.MigrateBulletInput{
		ID:     id,
		UserID: userID,
	})
	s.Nil(got)
	s.ErrorIs(gotErr, domain.ErrBulletNotOpenTask)
}

func (s *MigrateBulletUseCaseSuite) TestMigrateBullet_Success() {
	id := uuid.New()
	userID := uuid.New()
	dayAgo := time.Now().Add(-24 * time.Hour)

	bullet := &domain.Bullet{
		ID:        id,
		Kind:      domain.KindTask,
		Signifier: domain.SignifierOpen,
		Content:   "Test task",
		CreatedAt: dayAgo,
		UpdatedAt: dayAgo,
	}

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, id, userID).
		Return(bullet, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*domain.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.MigrateBulletInput{
		ID:     id,
		UserID: userID,
	})

	s.NoError(gotErr)
	s.Same(bullet, got)
	s.True(got.CreatedAt.After(dayAgo))
	s.True(got.UpdatedAt.After(dayAgo))
}

package usecase_test

import (
	"fmt"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CancelBulletUseCaseSuite struct {
	suite.Suite
	bulletRepoMock *mocks.MockBulletRepository
	uc             *usecase.CancelBullet
}

func TestCancelBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CancelBulletUseCaseSuite))
}

func (s *CancelBulletUseCaseSuite) SetupTest() {
	s.bulletRepoMock = mocks.NewMockBulletRepository(s.T())
	s.uc = usecase.NewCancelBullet(s.bulletRepoMock)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotFound() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CancelBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrNotFound)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_AlreadyCanceled() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierCancelled,
		}, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CancelBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.NoError(gotErr)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotOpen() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierCompleted,
		}, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CancelBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrBulletNotOpen)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_SaveError() {
	wantErr := fmt.Errorf("update failed")

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CancelBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_Success() {
	input := usecase.CancelBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, input.ID, input.UserID).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), input)
	s.NoError(gotErr)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

package cancelbullet_test

import (
	"fmt"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/cancelbullet"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/cancelbullet/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CancelBulletUseCaseSuite struct {
	suite.Suite
	bulletRepoMock *mocks.MockStore
	uc             *cancelbullet.Usecase
}

func TestCancelBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CancelBulletUseCaseSuite))
}

func (s *CancelBulletUseCaseSuite) SetupTest() {
	s.bulletRepoMock = mocks.NewMockStore(s.T())
	s.uc = cancelbullet.New(s.bulletRepoMock)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotFound() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), cancelbullet.Input{
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

	got, gotErr := s.uc.Execute(s.T().Context(), cancelbullet.Input{
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

	got, gotErr := s.uc.Execute(s.T().Context(), cancelbullet.Input{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrBulletNotOpen)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_UpdateError() {
	wantErr := fmt.Errorf("update failed")

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		UpdateSignifier(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), cancelbullet.Input{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_Success() {
	input := cancelbullet.Input{
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
		UpdateSignifier(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), input)
	s.NoError(gotErr)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_UpdateNotFound() {
	in := cancelbullet.Input{ID: uuid.New(), UserID: uuid.New()}

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, in.ID, in.UserID).
		Return(&entity.Bullet{Kind: entity.KindNote, Signifier: entity.SignifierOpen}, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		UpdateSignifier(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), in)
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrNotFound)
}

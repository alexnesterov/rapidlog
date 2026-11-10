package usecase_test

import (
	"fmt"
	"testing"

	mocksTransactor "github.com/alexnesterov/rapidlog-api/internal/pkg/transactor/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CancelBulletUseCaseSuite struct {
	suite.Suite
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocksTransactor.MockTransactor
	uc             *usecase.Usecase
}

func TestCancelBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CancelBulletUseCaseSuite))
}

func (s *CancelBulletUseCaseSuite) SetupTest() {
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocksTransactor.NewMockTransactor(s.T())
	s.uc = usecase.New(s.mockBulletRepo, s.mockTxMgr)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotFound() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.CancelBullet(s.T().Context(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrNotFound)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_AlreadyCanceled() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierCancelled,
		}, nil).
		Once()

	got, gotErr := s.uc.CancelBullet(s.T().Context(), uuid.New(), uuid.New())
	s.NoError(gotErr)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotOpen() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierCompleted,
		}, nil).
		Once()

	got, gotErr := s.uc.CancelBullet(s.T().Context(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrBulletNotOpen)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_UpdateError() {
	wantErr := fmt.Errorf("update failed")

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.CancelBullet(s.T().Context(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_Success() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.CancelBullet(s.T().Context(), uuid.New(), uuid.New())
	s.NoError(gotErr)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

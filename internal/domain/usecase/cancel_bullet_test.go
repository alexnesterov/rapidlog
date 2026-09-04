package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CancelBulletUseCaseSuite struct {
	suite.Suite
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocks.MockTransactionManager
	uc             port.BulletService
}

func TestCancelBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CancelBulletUseCaseSuite))
}

func (s *CancelBulletUseCaseSuite) SetupTest() {
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocks.NewMockTransactionManager(s.T())
	s.uc = usecase.NewBulletService(s.mockBulletRepo, s.mockTxMgr)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_Success() {
	s.mockBulletRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, err := s.uc.CancelBullet(context.Background(), uuid.New(), uuid.New())
	s.NoError(err)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotFound() {
	wantErr := port.ErrNotFound

	s.mockBulletRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, err := s.uc.CancelBullet(context.Background(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_AlreadyCanceled() {
	s.mockBulletRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierCancelled,
		}, nil).
		Once()

	got, err := s.uc.CancelBullet(context.Background(), uuid.New(), uuid.New())
	s.NoError(err)
	s.Equal(entity.SignifierCancelled, got.Signifier)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_NotOpen() {
	s.mockBulletRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierCompleted,
		}, nil).
		Once()

	got, err := s.uc.CancelBullet(context.Background(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, entity.ErrBulletNotOpen)
}

func (s *CancelBulletUseCaseSuite) TestCancelBullet_UpdateError() {
	wantErr := fmt.Errorf("update failed")

	s.mockBulletRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, err := s.uc.CancelBullet(context.Background(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

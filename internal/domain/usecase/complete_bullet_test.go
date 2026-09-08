package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CompleteBulletUseCaseSuite struct {
	suite.Suite
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocks.MockTransactionManager
	uc             port.BulletService
}

func TestCompleteBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CompleteBulletUseCaseSuite))
}

func (s *CompleteBulletUseCaseSuite) SetupTest() {
	s.mockBulletRepo = &mocks.MockBulletRepository{}
	s.mockTxMgr = &mocks.MockTransactionManager{}
	s.uc = usecase.NewBulletService(s.mockBulletRepo, s.mockTxMgr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_Success() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Type:      entity.BulletTask,
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, err := s.uc.CompleteBullet(context.Background(), uuid.New(), uuid.New())
	s.NoError(err)
	s.Equal(entity.SignifierCompleted, got.Signifier)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_NotFound() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, port.ErrNotFound).
		Once()

	got, err := s.uc.CompleteBullet(context.Background(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, port.ErrNotFound)
}

func (s *CompleteBulletUseCaseSuite) TestComleteBullet_AlreadyCompleted() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Type:      entity.BulletTask,
			Signifier: entity.SignifierCompleted,
		}, nil).
		Once()

	got, err := s.uc.CompleteBullet(context.Background(), uuid.New(), uuid.New())
	s.NoError(err)
	s.Equal(entity.SignifierCompleted, got.Signifier)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_ValidationError() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Type:      entity.BulletNote,
			Signifier: entity.SignifierCancelled,
		}, nil).
		Once()

	got, err := s.uc.CompleteBullet(context.Background(), uuid.New(), uuid.New())
	s.Nil(got)
	var validationErr *entity.ValidationError
	s.ErrorAs(err, &validationErr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_UpdateError() {
	wantErr := errors.New("update failed")

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Type:      entity.BulletTask,
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, err := s.uc.CompleteBullet(context.Background(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

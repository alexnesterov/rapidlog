package usecase_test

import (
	"errors"
	"testing"

	mocksTransactor "github.com/alexnesterov/rapidlog-api/internal/pkg/transactor/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CompleteBulletUseCaseSuite struct {
	suite.Suite
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocksTransactor.MockTransactor
	uc             *usecase.Usecase
}

func TestCompleteBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CompleteBulletUseCaseSuite))
}

func (s *CompleteBulletUseCaseSuite) SetupTest() {
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocksTransactor.NewMockTransactor(s.T())
	s.uc = usecase.New(s.mockBulletRepo, s.mockTxMgr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_NotFound() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.CompleteBullet(s.T().Context(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrNotFound)
}

func (s *CompleteBulletUseCaseSuite) TestComleteBullet_AlreadyCompleted() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindTask,
			Signifier: entity.SignifierCompleted,
		}, nil).
		Once()

	got, gotErr := s.uc.CompleteBullet(s.T().Context(), uuid.New(), uuid.New())
	s.NoError(gotErr)
	s.Equal(entity.SignifierCompleted, got.Signifier)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_ValidationError() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindNote,
			Signifier: entity.SignifierCancelled,
		}, nil).
		Once()

	got, gotErr := s.uc.CompleteBullet(s.T().Context(), uuid.New(), uuid.New())
	s.Nil(got)
	var validationErr *entity.ValidationError
	s.ErrorAs(gotErr, &validationErr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_UpdateError() {
	wantErr := errors.New("update failed")

	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindTask,
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.CompleteBullet(s.T().Context(), uuid.New(), uuid.New())
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_Success() {
	s.mockBulletRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindTask,
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Update(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.CompleteBullet(s.T().Context(), uuid.New(), uuid.New())
	s.NoError(gotErr)
	s.Equal(entity.SignifierCompleted, got.Signifier)
}

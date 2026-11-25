package usecase_test

import (
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CompleteBulletUseCaseSuite struct {
	suite.Suite
	bulletRepoMock *mocks.MockBulletRepository
	uc             *usecase.CompleteBullet
}

func TestCompleteBulletUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CompleteBulletUseCaseSuite))
}

func (s *CompleteBulletUseCaseSuite) SetupTest() {
	s.bulletRepoMock = mocks.NewMockBulletRepository(s.T())
	s.uc = usecase.NewCompleteBullet(s.bulletRepoMock)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_NotFound() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CompleteBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, entity.ErrNotFound)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_AlreadyCompleted() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindTask,
			Signifier: entity.SignifierCompleted,
		}, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CompleteBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.NoError(gotErr)
	s.Equal(entity.SignifierCompleted, got.Signifier)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_ValidationError() {
	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindNote,
			Signifier: entity.SignifierCancelled,
		}, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CompleteBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	var validationErr *entity.ValidationError
	s.ErrorAs(gotErr, &validationErr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_SaveError() {
	wantErr := errors.New("update failed")

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
		Return(&entity.Bullet{
			Kind:      entity.KindTask,
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.CompleteBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	})
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *CompleteBulletUseCaseSuite) TestCompleteBullet_Success() {
	input := usecase.CompleteBulletInput{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	s.bulletRepoMock.EXPECT().
		Get(mock.Anything, input.ID, input.UserID).
		Return(&entity.Bullet{
			Kind:      entity.KindTask,
			Signifier: entity.SignifierOpen,
		}, nil).
		Once()

	s.bulletRepoMock.EXPECT().
		Save(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), input)
	s.NoError(gotErr)
	s.Equal(entity.SignifierCompleted, got.Signifier)
}

package usecase_test

import (
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListBulletsUseCaseSuite struct {
	suite.Suite
	bulletRepoMock *mocks.MockBulletRepository
	uc             *usecase.ListBullets
}

func TestListBulletsUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ListBulletsUseCaseSuite))
}

func (s *ListBulletsUseCaseSuite) SetupTest() {
	s.bulletRepoMock = mocks.NewMockBulletRepository(s.T())
	s.uc = usecase.NewListBullets(s.bulletRepoMock)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_ListError() {
	wantErr := errors.New("list error")

	s.bulletRepoMock.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.ListBulletsInput{UserID: uuid.New()})
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Empty() {
	s.bulletRepoMock.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, nil).
		Once()

	userID := uuid.New()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.ListBulletsInput{UserID: userID})
	s.NoError(gotErr)
	s.Len(got, 0)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Success() {
	userID := uuid.New()
	want := []*domain.Bullet{{Content: "Заголовок", UserID: userID}}

	s.bulletRepoMock.EXPECT().
		List(mock.Anything, userID).
		Return(want, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), usecase.ListBulletsInput{UserID: userID})
	s.NoError(gotErr)
	s.Equal(want, got)
}

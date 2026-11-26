package listbullets_test

import (
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/listbullets"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/features/listbullets/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListBulletsUseCaseSuite struct {
	suite.Suite
	bulletRepoMock *mocks.MockStore
	uc             *listbullets.Usecase
}

func TestListBulletsUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ListBulletsUseCaseSuite))
}

func (s *ListBulletsUseCaseSuite) SetupTest() {
	s.bulletRepoMock = mocks.NewMockStore(s.T())
	s.uc = listbullets.New(s.bulletRepoMock)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_ListError() {
	wantErr := errors.New("list error")

	s.bulletRepoMock.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), listbullets.Input{UserID: uuid.New()})
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Empty() {
	s.bulletRepoMock.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, nil).
		Once()

	userID := uuid.New()

	got, gotErr := s.uc.Execute(s.T().Context(), listbullets.Input{UserID: userID})
	s.NoError(gotErr)
	s.Len(got, 0)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Success() {
	userID := uuid.New()
	want := []*entity.Bullet{{Content: "Заголовок", UserID: userID}}

	s.bulletRepoMock.EXPECT().
		List(mock.Anything, userID).
		Return(want, nil).
		Once()

	got, gotErr := s.uc.Execute(s.T().Context(), listbullets.Input{UserID: userID})
	s.NoError(gotErr)
	s.Equal(want, got)
}

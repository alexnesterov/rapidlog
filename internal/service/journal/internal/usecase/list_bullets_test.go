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

type ListBulletsUseCaseSuite struct {
	suite.Suite
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocksTransactor.MockTransactor
	uc             *usecase.Usecase
}

func TestListBulletsUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ListBulletsUseCaseSuite))
}

func (s *ListBulletsUseCaseSuite) SetupTest() {
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocksTransactor.NewMockTransactor(s.T())
	s.uc = usecase.New(s.mockBulletRepo, s.mockTxMgr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_ListError() {
	wantErr := errors.New("list error")

	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, gotErr := s.uc.ListBullets(s.T().Context(), uuid.New())
	s.Nil(got)
	s.ErrorIs(gotErr, wantErr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Empty() {
	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, nil).
		Once()

	userID := uuid.New()

	got, gotErr := s.uc.ListBullets(s.T().Context(), userID)
	s.NoError(gotErr)
	s.Len(got, 0)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Success() {
	userID := uuid.New()
	want := []*entity.Bullet{{Content: "Заголовок", UserID: userID}}

	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(want, nil).
		Once()

	got, gotErr := s.uc.ListBullets(s.T().Context(), userID)
	s.NoError(gotErr)
	s.Equal(want, got)
}

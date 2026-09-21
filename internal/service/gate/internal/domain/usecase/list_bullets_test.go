package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListBulletsUseCaseSuite struct {
	suite.Suite
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocks.MockTransactionManager
	uc             port.BulletService
}

func TestListBulletsUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ListBulletsUseCaseSuite))
}

func (s *ListBulletsUseCaseSuite) SetupTest() {
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocks.NewMockTransactionManager(s.T())
	s.uc = usecase.NewBulletService(s.mockBulletRepo, s.mockTxMgr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Success() {
	userID := uuid.New()
	want := []*entity.Bullet{{Content: "Заголовок", UserID: userID}}

	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(want, nil).
		Once()

	got, err := s.uc.ListBullets(context.Background(), userID)
	s.NoError(err)
	s.Equal(want, got)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_ListError() {
	wantErr := errors.New("list error")

	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, err := s.uc.ListBullets(context.Background(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_Empty() {
	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(nil).
		Times(3)

	s.mockTxMgr.EXPECT().WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		}).
		Once()

	userID := uuid.New()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	got, err := s.uc.ListBullets(context.Background(), userID)
	s.NoError(err)
	s.Len(got, 3)

	for _, b := range got {
		s.Equal(userID, b.UserID)
		s.Equal(yesterday, b.CreatedAt.Format("2006-01-02"))
	}
}

func (s *ListBulletsUseCaseSuite) TestListBullets_CreateError() {
	wantErr := errors.New("create bullets error")

	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, nil).
		Once()

	s.mockBulletRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*entity.Bullet")).
		Return(wantErr).
		Once()

	s.mockTxMgr.EXPECT().WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		}).
		Once()

	got, err := s.uc.ListBullets(context.Background(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

func (s *ListBulletsUseCaseSuite) TestListBullets_TransactionError() {
	wantErr := errors.New("transaction error")
	s.mockBulletRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, nil).
		Once()

	s.mockTxMgr.EXPECT().WithTransaction(mock.Anything, mock.Anything).
		Return(wantErr).
		Once()

	got, err := s.uc.ListBullets(context.Background(), uuid.New())
	s.Nil(got)
	s.ErrorIs(err, wantErr)
}

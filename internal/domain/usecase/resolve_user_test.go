package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ResolveUserUseCaseSuite struct {
	suite.Suite
	mockUserRepo   *mocks.MockUserRepository
	mockBulletRepo *mocks.MockBulletRepository
	mockTxMgr      *mocks.MockTransactionManager
	uc             port.UserService
}

func TestResolveUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ResolveUserUseCaseSuite))
}

func (s *ResolveUserUseCaseSuite) SetupTest() {
	s.mockUserRepo = mocks.NewMockUserRepository(s.T())
	s.mockBulletRepo = mocks.NewMockBulletRepository(s.T())
	s.mockTxMgr = mocks.NewMockTransactionManager(s.T())
	s.uc = usecase.NewUserService(s.mockUserRepo, s.mockBulletRepo, s.mockTxMgr)
}

func (s *ResolveUserUseCaseSuite) TestResolveUser_Exists() {
	id := uuid.New()

	user := &entity.User{ID: id, CreatedAt: time.Now()}

	s.mockUserRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(user, nil).
		Once()

	got, err := s.uc.ResolveUser(context.Background(), id)
	s.NoError(err)
	s.Equal(id, got)
}

func (s *ResolveUserUseCaseSuite) TestResolveUser_NotExists() {
	id := uuid.New()

	s.mockUserRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, port.ErrNotFound).
		Once()
	s.mockUserRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u != nil && u.ID != uuid.Nil })).
		Return(nil).
		Once()

	got, err := s.uc.ResolveUser(context.Background(), id)
	s.NoError(err)
	s.NotEqual(uuid.Nil, got)
	s.NotEqual(id, got)
}

func (s *ResolveUserUseCaseSuite) TestResolveUser_GetError() {
	id := uuid.New()
	wantErr := errors.New("get failed")

	s.mockUserRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, err := s.uc.ResolveUser(context.Background(), id)
	s.Equal(uuid.Nil, got)
	s.Error(err)
	s.ErrorIs(err, wantErr)
}

func (s *ResolveUserUseCaseSuite) TestResolveUser_CreateError() {
	id := uuid.New()
	wantErr := errors.New("create failed")

	s.mockUserRepo.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, port.ErrNotFound).
		Once()
	s.mockUserRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
		return u != nil && u.ID != uuid.Nil
	})).
		Return(wantErr).
		Once()

	got, err := s.uc.ResolveUser(context.Background(), id)
	s.Equal(uuid.Nil, got)
	s.Error(err)
	s.ErrorIs(err, wantErr)
}

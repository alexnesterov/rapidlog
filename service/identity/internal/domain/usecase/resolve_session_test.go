package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ResolveSessionUseCaseSuite struct {
	suite.Suite
	mockUserRepo *mocks.MockUserRepository
	uc           port.IdentityService
}

func TestResolveSessionUseCaseSuite(t *testing.T) {
	suite.Run(t, &ResolveSessionUseCaseSuite{})
}

func (s *ResolveSessionUseCaseSuite) SetupTest() {
	s.mockUserRepo = mocks.NewMockUserRepository(s.T())
	s.uc = usecase.NewIdentityService(s.mockUserRepo)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_Exists() {
	user := entity.NewUser()

	s.mockUserRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(user, nil).
		Once()

	got, err := s.uc.ResolveSession(context.Background(), user.ID)
	s.NoError(err)
	s.Equal(user.ID, got)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_NotExists() {
	user := entity.NewUser()

	s.mockUserRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, port.ErrNotFound).
		Once()

	s.mockUserRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u != nil && u.ID != uuid.Nil })).
		Return(nil).
		Once()

	got, err := s.uc.ResolveSession(context.Background(), user.ID)
	s.NoError(err)
	s.NotEqual(user.ID, got)
	s.NotEqual(uuid.Nil, got)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_GetError() {
	user := entity.NewUser()
	wantErr := errors.New("get error")

	s.mockUserRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, err := s.uc.ResolveSession(context.Background(), user.ID)
	s.ErrorIs(err, wantErr)
	s.Equal(uuid.Nil, got)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_CreateError() {
	user := entity.NewUser()
	wantErr := errors.New("create error")

	s.mockUserRepo.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, port.ErrNotFound).
		Once()

	s.mockUserRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u != nil && u.ID != uuid.Nil })).
		Return(wantErr).
		Once()

	got, err := s.uc.ResolveSession(context.Background(), user.ID)
	s.ErrorIs(err, wantErr)
	s.Equal(uuid.Nil, got)
}

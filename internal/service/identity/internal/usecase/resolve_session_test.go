package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/usecase"
	"github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ResolveSessionUseCaseSuite struct {
	suite.Suite
	userRepoMock *mocks.MockUserRepository
	uc           *usecase.Usecase
}

func TestResolveSessionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(ResolveSessionUseCaseSuite))
}

func (s *ResolveSessionUseCaseSuite) SetupTest() {
	s.userRepoMock = mocks.NewMockUserRepository(s.T())
	s.uc = usecase.New(s.userRepoMock)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_Exists() {
	user := entity.NewUser()

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(user, nil).
		Once()

	got, gotErr := s.uc.ResolveSession(context.Background(), user.ID)
	s.NoError(gotErr)
	s.Equal(user.ID, got)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_NotExists() {
	user := entity.NewUser()

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	s.userRepoMock.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u != nil && u.ID != uuid.Nil })).
		Return(nil).
		Once()

	got, gotErr := s.uc.ResolveSession(context.Background(), user.ID)
	s.NoError(gotErr)
	s.NotEqual(user.ID, got)
	s.NotEqual(uuid.Nil, got)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_GetError() {
	user := entity.NewUser()
	wantErr := errors.New("get error")

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, gotErr := s.uc.ResolveSession(context.Background(), user.ID)
	s.ErrorIs(gotErr, wantErr)
	s.Equal(uuid.Nil, got)
}

func (s *ResolveSessionUseCaseSuite) TestResolveSession_CreateError() {
	user := entity.NewUser()
	wantErr := errors.New("create error")

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	s.userRepoMock.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u != nil && u.ID != uuid.Nil })).
		Return(wantErr).
		Once()

	got, gotErr := s.uc.ResolveSession(context.Background(), user.ID)
	s.ErrorIs(gotErr, wantErr)
	s.Equal(uuid.Nil, got)
}

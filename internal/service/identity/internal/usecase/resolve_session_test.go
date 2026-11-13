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

type ResolveSessionSuite struct {
	suite.Suite
	userRepoMock *mocks.MockUserRepository
	uc           *usecase.ResolveSession
}

func TestResolveSessionSuite(t *testing.T) {
	suite.Run(t, new(ResolveSessionSuite))
}

func (s *ResolveSessionSuite) SetupTest() {
	s.userRepoMock = mocks.NewMockUserRepository(s.T())
	s.uc = usecase.NewResolveSession(s.userRepoMock)
}

func (s *ResolveSessionSuite) TestResolveSession_Exists() {
	user := entity.NewUser()

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(user, nil).
		Once()

	got, gotErr := s.uc.Execute(context.Background(), user.ID)
	s.NoError(gotErr)
	s.Equal(user.ID, got)
}

func (s *ResolveSessionSuite) TestResolveSession_NotExists() {
	user := entity.NewUser()

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, entity.ErrNotFound).
		Once()

	s.userRepoMock.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(u *entity.User) bool { return u != nil && u.ID != uuid.Nil })).
		Return(nil).
		Once()

	got, gotErr := s.uc.Execute(context.Background(), user.ID)
	s.NoError(gotErr)
	s.NotEqual(user.ID, got)
	s.NotEqual(uuid.Nil, got)
}

func (s *ResolveSessionSuite) TestResolveSession_GetError() {
	user := entity.NewUser()
	wantErr := errors.New("get error")

	s.userRepoMock.EXPECT().
		Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, wantErr).
		Once()

	got, gotErr := s.uc.Execute(context.Background(), user.ID)
	s.ErrorIs(gotErr, wantErr)
	s.Equal(uuid.Nil, got)
}

func (s *ResolveSessionSuite) TestResolveSession_CreateError() {
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

	got, gotErr := s.uc.Execute(s.T().Context(), user.ID)
	s.ErrorIs(gotErr, wantErr)
	s.Equal(uuid.Nil, got)
}

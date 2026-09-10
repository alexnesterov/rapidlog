package usecase_test

import (
	"context"
	"testing"

	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/port"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/port/mocks"
	"github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/usecase"
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

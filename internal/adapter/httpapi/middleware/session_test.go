package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const testCookieName = "session_id"

type SessionMiddlewareSuite struct {
	suite.Suite
	mockUserService *mocks.MockUserService
}

func TestSessionMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(SessionMiddlewareSuite))
}

func (s *SessionMiddlewareSuite) SetupTest() {
	s.mockUserService = mocks.NewMockUserService(s.T())
}

func (s *SessionMiddlewareSuite) TestSession_NoCookie_SetsNewCookie() {
	newID := uuid.New()

	s.mockUserService.EXPECT().ResolveUser(mock.Anything, uuid.Nil).
		Return(newID, nil).
		Once()

	var gotID uuid.UUID
	var gotOK bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID, gotOK = httpapi.UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := Session(s.mockUserService, testCookieName, 365*24*time.Hour, true)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	s.True(gotOK)
	s.Equal(newID, gotID)

	cookies := res.Result().Cookies()
	s.Len(cookies, 1)
	s.Equal(testCookieName, cookies[0].Name)
	s.Equal(newID.String(), cookies[0].Value)
	s.Equal("/", cookies[0].Path)
	s.True(cookies[0].HttpOnly)
	s.True(cookies[0].Secure)
	s.Equal(int((365 * 24 * time.Hour).Seconds()), cookies[0].MaxAge)
	s.Equal(http.SameSiteLaxMode, cookies[0].SameSite)
}

func (s *SessionMiddlewareSuite) TestSession_ValidCookie_NoNewCookie() {
	id := uuid.New()

	s.mockUserService.EXPECT().ResolveUser(mock.Anything, id).
		Return(id, nil).
		Once()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Session(s.mockUserService, testCookieName, 365*24*time.Hour, true)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: id.String()})
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	s.Empty(res.Result().Cookies())
}

func (s *SessionMiddlewareSuite) TestSession_StaleCookie_SetsNewCookie() {
	staleID := uuid.New()
	newID := uuid.New()

	s.mockUserService.EXPECT().ResolveUser(mock.Anything, staleID).
		Return(newID, nil).
		Once()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Session(s.mockUserService, testCookieName, 365*24*time.Hour, true)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: staleID.String()})
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	cookies := res.Result().Cookies()
	s.Len(cookies, 1)
	s.Equal(newID.String(), cookies[0].Value)
}

func (s *SessionMiddlewareSuite) TestSession_HealthPath_Bypassed() {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := Session(s.mockUserService, testCookieName, 365*24*time.Hour, true)(next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	s.True(nextCalled)
	s.mockUserService.AssertNotCalled(s.T(), "ResolveUser", mock.Anything, mock.Anything)
}

func (s *SessionMiddlewareSuite) TestSession_ResolveError_Returns500() {
	wantErr := errors.New("resolve failed")

	s.mockUserService.EXPECT().ResolveUser(mock.Anything, uuid.Nil).
		Return(uuid.Nil, wantErr).
		Once()

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	handler := Session(s.mockUserService, testCookieName, 365*24*time.Hour, true)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	s.False(nextCalled)
	s.Equal(http.StatusInternalServerError, res.Code)
}

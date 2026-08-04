package httpapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name       string
		setupMock  func(*mocks.MockPinger)
		wantStatus int
		wantBody   string
	}{
		{
			name: "healthy",
			setupMock: func(m *mocks.MockPinger) {
				m.EXPECT().Ping(mock.Anything).Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok","db":"ok"}`,
		},
		{
			name: "database error",
			setupMock: func(m *mocks.MockPinger) {
				m.EXPECT().Ping(mock.Anything).Return(errors.New("db error"))
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `{"status":"error","db":"error"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockPinger := mocks.NewMockPinger(t)
			tc.setupMock(mockPinger)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			res := httptest.NewRecorder()

			handler := NewHealthHandler(mockPinger)
			handler(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
			assert.JSONEq(t, tc.wantBody, res.Body.String())
		})
	}
}

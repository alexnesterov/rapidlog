package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPinger struct {
	pingErr error
}

func (m *mockPinger) Ping(ctx context.Context) error {
	return m.pingErr
}

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name       string
		pinger     Pinger
		wantStatus int
		wantBody   string
	}{
		{
			name:       "healthy",
			pinger:     &mockPinger{pingErr: nil},
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok","db":"ok"}`,
		},
		{
			name:       "database error",
			pinger:     &mockPinger{pingErr: errors.New("db error")},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `{"status":"error","db":"error"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			res := httptest.NewRecorder()

			handler := NewHealthHandler(tc.pinger)
			handler(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
			assert.JSONEq(t, tc.wantBody, res.Body.String())
		})
	}
}

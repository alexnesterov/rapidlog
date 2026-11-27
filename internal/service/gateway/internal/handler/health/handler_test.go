package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/handler/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "healthy",
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			res := httptest.NewRecorder()

			h := health.NewHandler()
			h.Handle(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
			assert.JSONEq(t, tc.wantBody, res.Body.String())
		})
	}
}

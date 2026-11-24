package middleware

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/kratos"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type resolverFunc func(ctx context.Context, cookie, sessionToken string) (uuid.UUID, error)

func (f resolverFunc) Whoami(ctx context.Context, cookie, sessionToken string) (uuid.UUID, error) {
	return f(ctx, cookie, sessionToken)
}

func TestAuth(t *testing.T) {
	userID := uuid.New()

	cases := []struct {
		name       string
		resolveID  uuid.UUID
		resolveErr error
		wantStatus int
		wantBody   string
		wantNext   bool
	}{
		{
			name:       "active session",
			resolveID:  userID,
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
		{
			name:       "no session",
			resolveErr: kratos.ErrUnauthenticated,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"data":null,"error":{"code":401,"message":"unauthorized"}}`,
		},
		{
			name:       "kratos unavailable",
			resolveErr: errors.New("connection refused"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"data":null,"error":{"code":500,"message":"internal server error"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotCookie, gotToken string
			resolver := resolverFunc(func(_ context.Context, cookie, sessionToken string) (uuid.UUID, error) {
				gotCookie, gotToken = cookie, sessionToken
				return tc.resolveID, tc.resolveErr
			})

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				id, ok := httpapi.UserIDFromContext(r.Context())
				assert.True(t, ok)
				assert.Equal(t, userID, id)
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/api/bullets", nil)
			req.Header.Set("Cookie", "ory_kratos_session=abc")
			req.Header.Set("X-Session-Token", "tok")
			res := httptest.NewRecorder()

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			Auth(resolver, logger)(next).ServeHTTP(res, req)

			require.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, tc.wantNext, nextCalled)
			assert.Equal(t, "ory_kratos_session=abc", gotCookie)
			assert.Equal(t, "tok", gotToken)
			if tc.wantBody != "" {
				assert.JSONEq(t, tc.wantBody, res.Body.String())
			}
		})
	}
}

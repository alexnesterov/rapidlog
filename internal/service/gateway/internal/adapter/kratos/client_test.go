package kratos

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Whoami(t *testing.T) {
	identityID := uuid.New()

	cases := []struct {
		name         string
		cookie       string
		sessionToken string
		status       int
		body         string
		wantID       uuid.UUID
		wantErr      error
		wantAnyErr   bool
		wantNoCall   bool
	}{
		{
			name:   "active session by cookie",
			cookie: "ory_kratos_session=abc",
			status: http.StatusOK,
			body:   `{"active":true,"identity":{"id":"` + identityID.String() + `"}}`,
			wantID: identityID,
		},
		{
			name:         "active session by token",
			sessionToken: "tok",
			status:       http.StatusOK,
			body:         `{"active":true,"identity":{"id":"` + identityID.String() + `"}}`,
			wantID:       identityID,
		},
		{
			name:       "no credentials skips kratos",
			wantErr:    ErrUnauthenticated,
			wantNoCall: true,
		},
		{
			name:    "unauthorized",
			cookie:  "ory_kratos_session=expired",
			status:  http.StatusUnauthorized,
			body:    `{"error":{"code":401}}`,
			wantErr: ErrUnauthenticated,
		},
		{
			name:    "forbidden (aal too low)",
			cookie:  "ory_kratos_session=abc",
			status:  http.StatusForbidden,
			body:    `{"error":{"code":403}}`,
			wantErr: ErrUnauthenticated,
		},
		{
			name:    "inactive session",
			cookie:  "ory_kratos_session=abc",
			status:  http.StatusOK,
			body:    `{"active":false,"identity":{"id":"` + identityID.String() + `"}}`,
			wantErr: ErrUnauthenticated,
		},
		{
			name:       "kratos failure",
			cookie:     "ory_kratos_session=abc",
			status:     http.StatusInternalServerError,
			body:       `{}`,
			wantAnyErr: true,
		},
		{
			name:       "malformed body",
			cookie:     "ory_kratos_session=abc",
			status:     http.StatusOK,
			body:       `not json`,
			wantAnyErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				assert.Equal(t, "/sessions/whoami", r.URL.Path)
				assert.Equal(t, tc.cookie, r.Header.Get("Cookie"))
				assert.Equal(t, tc.sessionToken, r.Header.Get("X-Session-Token"))
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			client, err := NewClient(srv.URL)
			require.NoError(t, err)

			id, err := client.Whoami(context.Background(), tc.cookie, tc.sessionToken)

			assert.Equal(t, !tc.wantNoCall, called)
			switch {
			case tc.wantErr != nil:
				require.ErrorIs(t, err, tc.wantErr)
			case tc.wantAnyErr:
				require.Error(t, err)
				assert.NotErrorIs(t, err, ErrUnauthenticated)
			default:
				require.NoError(t, err)
				assert.Equal(t, tc.wantID, id)
			}
		})
	}
}

func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient("kratos:4433")
	require.Error(t, err)
}

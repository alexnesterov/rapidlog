package httpapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListBullets(t *testing.T) {
	userID := uuid.New()

	cases := []struct {
		name string
		wantStatus int
		wantData []httpapi.BulletDayGroup
	} {
		{
			name: "success",
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/bullets", nil)
			res := httptest.NewRecorder()

			httpapi.ListBullets(res, req.WithContext(httpapi.WithUserID(context.Background(), userID)))

			assert.Equal(t, tc.wantStatus, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		})
	}
}
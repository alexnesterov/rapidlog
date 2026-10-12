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

	req := httptest.NewRequest(http.MethodGet, "/api/bullets", nil)
	res := httptest.NewRecorder()

	httpapi.ListBullets(res, req.WithContext(httpapi.WithUserID(context.Background(), userID)))

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}
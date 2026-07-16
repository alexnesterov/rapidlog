package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	handler := http.HandlerFunc(Health)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)

	handler.ServeHTTP(res, req)

	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, `{"status": "ok"}`, res.Body.String())
}

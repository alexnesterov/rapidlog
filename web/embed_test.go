package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestSPAHandler(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html":       {Data: []byte("<!doctype html>app")},
		"assets/app-1.js":  {Data: []byte("console.log(1)")},
		"assets/app-1.css": {Data: []byte("body{}")},
	}

	cases := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "root", path: "/", wantStatus: http.StatusOK, wantBody: "<!doctype html>app"},
		{name: "asset", path: "/assets/app-1.js", wantStatus: http.StatusOK, wantBody: "console.log(1)"},
		{name: "client route", path: "/login", wantStatus: http.StatusOK, wantBody: "<!doctype html>app"},
		{name: "client route with query", path: "/registration?flow=abc", wantStatus: http.StatusOK, wantBody: "<!doctype html>app"},
		{name: "nested client route", path: "/journal/2026-11-17", wantStatus: http.StatusOK, wantBody: "<!doctype html>app"},
		{name: "missing asset stays 404", path: "/assets/app-0.js", wantStatus: http.StatusNotFound},
		{name: "missing file with extension stays 404", path: "/favicon.ico", wantStatus: http.StatusNotFound},
	}

	h := spaHandler(fsys)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			h.ServeHTTP(res, httptest.NewRequest(http.MethodGet, tc.path, nil))

			assert.Equal(t, tc.wantStatus, res.Code)
			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, res.Body.String())
			}
		})
	}
}

package kratos

import (
	"net/http"
	"net/http/httputil"
)

// NewProxy проксирует публичный API Kratos (self-service flows, whoami, logout) для браузера.
// Монтируется под префиксом (/auth/), который снимает вызывающий через http.StripPrefix:
// так cookie сессии Kratos живёт на том же origin, что и SPA, и CORS не нужен.
func NewProxy(publicURL string) (http.Handler, error) {
	target, err := parsePublicURL(publicURL)
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			pr.SetXForwarded()
		},
	}, nil
}

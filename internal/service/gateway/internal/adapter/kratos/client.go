// Package kratos — адаптер к публичному API Ory Kratos: проверка сессии и прокси
// self-service эндпоинтов для браузера.
package kratos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// ErrUnauthenticated — у запроса нет действующей сессии Kratos.
var ErrUnauthenticated = errors.New("kratos: unauthenticated")

const whoamiTimeout = 2 * time.Second

type Client struct {
	whoamiURL string
	http      *http.Client
}

func NewClient(publicURL string) (*Client, error) {
	u, err := parsePublicURL(publicURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		whoamiURL: u.JoinPath("sessions", "whoami").String(),
		http:      &http.Client{Timeout: whoamiTimeout},
	}, nil
}

type whoamiResponse struct {
	Active   bool `json:"active"`
	Identity struct {
		ID uuid.UUID `json:"id"`
	} `json:"identity"`
}

// Whoami возвращает ID identity Kratos для сессии из cookie браузера или токена
// API-клиента (заголовок X-Session-Token).
func (c *Client) Whoami(ctx context.Context, cookie, sessionToken string) (uuid.UUID, error) {
	if cookie == "" && sessionToken == "" {
		return uuid.Nil, ErrUnauthenticated
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.whoamiURL, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("kratos: whoami: new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if sessionToken != "" {
		req.Header.Set("X-Session-Token", sessionToken)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("kratos: whoami: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized, http.StatusForbidden:
		return uuid.Nil, ErrUnauthenticated
	default:
		return uuid.Nil, fmt.Errorf("kratos: whoami: unexpected status %d", resp.StatusCode)
	}

	var body whoamiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return uuid.Nil, fmt.Errorf("kratos: whoami: decode: %w", err)
	}
	if !body.Active || body.Identity.ID == uuid.Nil {
		return uuid.Nil, ErrUnauthenticated
	}

	return body.Identity.ID, nil
}

func parsePublicURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("kratos: parse public url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("kratos: public url %q must be absolute", raw)
	}
	return u, nil
}

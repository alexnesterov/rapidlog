package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/kratos"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"github.com/google/uuid"
)

// SessionResolver проверяет сессию Kratos по cookie браузера или токену API-клиента.
type SessionResolver interface {
	Whoami(ctx context.Context, cookie, sessionToken string) (uuid.UUID, error)
}

// Auth пускает дальше только запросы с действующей сессией Kratos и кладёт ID её identity
// в контекст как user_id. Без сессии — 401, недоступный Kratos — 500.
func Auth(resolver SessionResolver, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := resolver.Whoami(r.Context(), r.Header.Get("Cookie"), r.Header.Get("X-Session-Token"))
			if errors.Is(err, kratos.ErrUnauthenticated) {
				respond.Error(w, http.StatusUnauthorized, httpapi.ErrUnauthorized.Error())
				return
			}
			if err != nil {
				logger.Error("failed to resolve kratos session", "error", err)
				respond.Error(w, http.StatusInternalServerError, httpapi.ErrInternal.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(httpapi.WithUserID(r.Context(), userID)))
		})
	}
}

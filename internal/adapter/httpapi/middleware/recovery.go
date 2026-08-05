package middleware

import (
	"log/slog"
	"net/http"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
					)

					httpapi.RespondError(w, http.StatusInternalServerError, "internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

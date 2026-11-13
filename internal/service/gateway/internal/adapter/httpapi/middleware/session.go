package middleware

import (
	"log/slog"
	"net/http"
	"time"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/pkg/respond"
	"github.com/google/uuid"
)

func Session(identityClient identityv1.IdentityServiceClient, cookieName string, maxAge time.Duration, secure bool, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			var id uuid.UUID
			if c, err := r.Cookie(cookieName); err == nil {
				parsed, perr := uuid.Parse(c.Value)
				if perr == nil {
					id = parsed
				}
			}

			resp, err := identityClient.ResolveSession(r.Context(), &identityv1.ResolveSessionRequest{SessionId: id.String()})
			if err != nil {
				logger.Error("failed to resolve session", "error", err)
				respond.Error(w, http.StatusInternalServerError, httpapi.ErrInternal.Error())
				return
			}

			resolved, err := uuid.Parse(resp.GetUserId())
			if err != nil {
				logger.Error("failed to parse user id", "error", err)
				respond.Error(w, http.StatusInternalServerError, httpapi.ErrInternal.Error())
				return
			}

			ctx := httpapi.WithUserID(r.Context(), resolved)
			r = r.WithContext(ctx)

			if resolved != id {
				http.SetCookie(w, &http.Cookie{
					Name:     cookieName,
					Value:    resolved.String(),
					Path:     "/",
					MaxAge:   int(maxAge.Seconds()),
					HttpOnly: true,
					Secure:   secure,
					SameSite: http.SameSiteLaxMode,
				})
			}

			next.ServeHTTP(w, r)
		})
	}
}

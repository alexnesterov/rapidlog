package middleware

import (
	"context"
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
			var id uuid.UUID
			if c, err := r.Cookie(cookieName); err == nil {
				if parsed, perr := uuid.Parse(c.Value); perr == nil {
					id = parsed
				}
			}

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			resp, err := identityClient.ResolveSession(ctx, &identityv1.ResolveSessionRequest{SessionId: id.String()})
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

			r = r.WithContext(httpapi.WithUserID(r.Context(), resolved))

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

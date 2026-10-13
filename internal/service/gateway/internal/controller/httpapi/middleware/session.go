package middleware

import (
	"net/http"
	"time"

	identityv1 "github.com/alexnesterov/rapidlog-api/api/gen/identity/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/gateway/internal/controller/httpapi"
	"github.com/google/uuid"
)

func Session(identityClient identityv1.IdentityServiceClient, cookieName string, maxAge time.Duration, secure bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			var id uuid.UUID
			if c, err := r.Cookie(cookieName); err == nil {
				id, _ = uuid.Parse(c.Value)
			}

			resp, err := identityClient.ResolveSession(r.Context(), &identityv1.ResolveSessionRequest{SessionId: id.String()})
			if err != nil {
				httpapi.RespondError(w, http.StatusInternalServerError, httpapi.ErrInternal.Error())
				return
			}

			resolved := uuid.MustParse(resp.GetUserId())

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

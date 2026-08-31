package middleware

import (
	"net/http"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/adapter/httpapi"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

func Session(userService port.UserService, cookieName string, maxAge time.Duration, secure bool) func(http.Handler) http.Handler {
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

			resolved, err := userService.ResolveUser(r.Context(), id)
			if err != nil {
				httpapi.RespondError(w, http.StatusInternalServerError, httpapi.ErrInternal.Error())
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

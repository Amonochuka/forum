package middleware

import (
	"context"
	"forum/internal/session"
	"net/http"
)

func RequireAuth(sessionsService *session.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, err := sessionsService.ValidateSession(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
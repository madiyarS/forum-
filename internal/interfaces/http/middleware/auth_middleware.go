package middleware

import (
	"context"
	"net/http"
	"time"

	"forum/internal/usecases/auth"
)

// AuthMiddleware authenticates requests
func AuthMiddleware(service *auth.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sessionID, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			session, err := service.GetSessionRepo().FindByID(sessionID.Value)
			if err != nil || session == nil || session.ExpiresAt.Before(time.Now()) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", session.UserID)
			ctx = context.WithValue(ctx, "isLoggedIn", true)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
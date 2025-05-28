package middleware

import (
	"net/http"
	"time"

	"forum/models"
)

func AuthMiddleware(app *models.App) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sessionID, err := r.Cookie("session_id")
			if err != nil || sessionID == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			var userID int
			var expiresAt time.Time
			err = app.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE id = ?", sessionID.Value).Scan(&userID, &expiresAt)
			if err != nil || time.Now().After(expiresAt) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next(w, r)
		}
	}
}

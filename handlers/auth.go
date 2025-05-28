package handlers

import (
	"net/http"
	"time"

	"forum/models"
	"forum/utils"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")

			if username == "" || email == "" || password == "" {
				utils.RenderError(app, w, "All fields are required", http.StatusBadRequest)
				return
			}

			var count int
			err := app.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ?", email, username).Scan(&count)
			if err != nil || count > 0 {
				utils.RenderError(app, w, "Email or username already exists", http.StatusBadRequest)
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			_, err = app.DB.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, email, string(hash))
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		app.Templates.ExecuteTemplate(w, "register.html", nil)
	}
}

func Login(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			if email == "" || password == "" {
				utils.RenderError(app, w, "Email and password are required", http.StatusBadRequest)
				return
			}

			var user struct {
				ID           int
				PasswordHash string
			}
			err := app.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email).Scan(&user.ID, &user.PasswordHash)
			if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
				utils.RenderError(app, w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			_, err = app.DB.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			sessionID, _ := uuid.NewV4()
			expiresAt := time.Now().Add(24 * time.Hour)
			_, err = app.DB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID.String(), user.ID, expiresAt)
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID.String(),
				Expires:  expiresAt,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		app.Templates.ExecuteTemplate(w, "login.html", nil)
	}
}

func Logout(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err == nil {
			_, err = app.DB.Exec("DELETE FROM sessions WHERE id = ?", sessionID.Value)
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
			Secure:   true,
		})

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

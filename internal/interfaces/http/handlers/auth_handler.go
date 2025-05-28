package handlers

import (
	"html/template"
	"net/http"
	"time"

	"forum/internal/usecases/auth"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	service *auth.Service
	tmpl    *template.Template
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *auth.Service, tmpl *template.Template) *AuthHandler {
	return &AuthHandler{
		service: service,
		tmpl:    tmpl,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := h.service.Register(username, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	h.tmpl.ExecuteTemplate(w, "register.html", nil)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		session, err := h.service.Login(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.tmpl.ExecuteTemplate(w, "login.html", nil)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err == nil {
		h.service.Logout(sessionID.Value)
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
package handlers

import (
	"net/http"
	"strings"
	"time"

	"forum/models"
	"forum/utils"
)

func Comment(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RenderError(app, w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionID, _ := r.Cookie("session_id")
		var userID int
		err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
		if err != nil {
			utils.RenderError(app, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postID := r.FormValue("post_id")
		content := strings.TrimSpace(r.FormValue("content"))

		if content == "" {
			utils.RenderError(app, w, "Comment content is required", http.StatusBadRequest)
			return
		}

		_, err = app.DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
		if err != nil {
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
	}
}

func Like(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.RenderError(app, w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionID, _ := r.Cookie("session_id")
		var userID int
		err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
		if err != nil {
			utils.RenderError(app, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postID := r.FormValue("post_id")
		commentID := r.FormValue("comment_id")
		isLike := r.FormValue("is_like") == "true"

		tx, err := app.DB.Begin()
		if err != nil {
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if postID != "" {
			_, err = tx.Exec("INSERT OR REPLACE INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, isLike)
		} else if commentID != "" {
			_, err = tx.Exec("INSERT OR REPLACE INTO likes (user_id, comment_id, is_like) VALUES (?, ?, ?)", userID, commentID, isLike)
		} else {
			tx.Rollback()
			utils.RenderError(app, w, "Invalid request", http.StatusBadRequest)
			return
		}

		if err != nil {
			tx.Rollback()
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

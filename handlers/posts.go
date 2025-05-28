package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/models"
	"forum/utils"
)

func Index(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			utils.RenderError(app, w, "Page not found", http.StatusNotFound)
			return
		}

		rows, err := app.DB.Query("SELECT id, user_id, title, content, created_at FROM posts ORDER BY created_at DESC")
		if err != nil {
			log.Printf("Failed to query posts: %v", err)
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var p models.Post
			if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
				log.Printf("Failed to scan post: %v", err)
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			posts = append(posts, p)
		}
		log.Printf("Fetched %d posts", len(posts))

		rows, err = app.DB.Query("SELECT id, name FROM categories")
		if err != nil {
			log.Printf("Failed to query categories: %v", err)
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var categories []struct {
			ID   int
			Name string
		}
		for rows.Next() {
			var c struct {
				ID   int
				Name string
			}
			rows.Scan(&c.ID, &c.Name)
			categories = append(categories, c)
		}
		log.Printf("Fetched %d categories", len(categories))

		sessionID, _ := r.Cookie("session_id")
		isLoggedIn := false
		if sessionID != nil {
			var userID int
			err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
			if err == nil {
				isLoggedIn = true
			}
		}

		err = app.Templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"Title":      "Home",
			"Posts":      posts,
			"Categories": categories,
			"IsLoggedIn": isLoggedIn,
		})
		if err != nil {
			log.Printf("Failed to render index.html: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func NewPost(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sessionID, _ := r.Cookie("session_id")
			var userID int
			err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
			if err != nil {
				utils.RenderError(app, w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			title := strings.TrimSpace(r.FormValue("title"))
			content := strings.TrimSpace(r.FormValue("content"))
			categories := r.Form["categories"]

			if title == "" || content == "" {
				utils.RenderError(app, w, "Title and content are required", http.StatusBadRequest)
				return
			}

			tx, err := app.DB.Begin()
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			result, err := tx.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
			if err != nil {
				tx.Rollback()
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			postID, _ := result.LastInsertId()

			for _, cat := range categories {
				_, err = tx.Exec("INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, cat)
				if err != nil {
					tx.Rollback()
					utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			if err := tx.Commit(); err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/posts/%d", postID), http.StatusSeeOther)
			return
		}

		rows, err := app.DB.Query("SELECT id, name FROM categories")
		if err != nil {
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var categories []struct {
			ID   int
			Name string
		}
		for rows.Next() {
			var c struct {
				ID   int
				Name string
			}
			rows.Scan(&c.ID, &c.Name)
			categories = append(categories, c)
		}

		app.Templates.ExecuteTemplate(w, "new_post.html", map[string]interface{}{
			"Title":      "New Post",
			"Categories": categories,
		})
	}
}

func Find(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			utils.RenderError(app, w, "Page not found", http.StatusNotFound)
			return
		}

		postID, err := strconv.Atoi(parts[2])
		if err != nil {
			utils.RenderError(app, w, "Page not found", http.StatusNotFound)
			return
		}

		var post models.Post
		err = app.DB.QueryRow(`
            SELECT id, user_id, title, content, created_at
            FROM posts
            WHERE id = ?`, postID).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			utils.RenderError(app, w, "Page not found", http.StatusNotFound)
			return
		}

		rows, err := app.DB.Query(`
            SELECT c.id, c.post_id, c.user_id, c.content, c.created_at
            FROM comments c
            WHERE c.post_id = ?
            ORDER BY c.created_at ASC`, postID)
		if err != nil {
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []struct {
			models.Comment
			LikeCount    int
			DislikeCount int
			HasLiked     bool
			HasDisliked  bool
		}
		for rows.Next() {
			var c struct {
				models.Comment
				LikeCount    int
				DislikeCount int
				HasLiked     bool
				HasDisliked  bool
			}
			if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			err = app.DB.QueryRow(`
                SELECT 
                    COUNT(CASE WHEN is_like = 1 THEN 1 END),
                    COUNT(CASE WHEN is_like = 0 THEN 1 END)
                FROM likes
                WHERE comment_id = ?`, c.ID).Scan(&c.LikeCount, &c.DislikeCount)
			if err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			comments = append(comments, c)
		}

		var likeCount, dislikeCount int
		err = app.DB.QueryRow(`
            SELECT 
                COUNT(CASE WHEN is_like = 1 THEN 1 END),
                COUNT(CASE WHEN is_like = 0 THEN 1 END)
            FROM likes
            WHERE post_id = ?`, postID).Scan(&likeCount, &dislikeCount)
		if err != nil {
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		sessionID, _ := r.Cookie("session_id")
		isLoggedIn := false
		var userID int
		if sessionID != nil {
			err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
			if err == nil {
				isLoggedIn = true
			}
		}

		hasLiked := false
		hasDisliked := false
		if isLoggedIn {
			var isLike bool
			err := app.DB.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&isLike)
			if err == nil {
				if isLike {
					hasLiked = true
				} else {
					hasDisliked = true
				}
			}
		}

		for i, c := range comments {
			if isLoggedIn {
				var isLike bool
				err := app.DB.QueryRow("SELECT is_like FROM likes WHERE user_id = ? AND comment_id = ?", userID, c.ID).Scan(&isLike)
				if err == nil {
					if isLike {
						comments[i].HasLiked = true
					} else {
						comments[i].HasDisliked = true
					}
				}
			}
		}

		app.Templates.ExecuteTemplate(w, "post.html", map[string]interface{}{
			"Title":        "Post",
			"Post":         post,
			"Comments":     comments,
			"LikeCount":    likeCount,
			"DislikeCount": dislikeCount,
			"IsLoggedIn":   isLoggedIn,
			"HasLiked":     hasLiked,
			"HasDisliked":  hasDisliked,
		})
	}
}

func Filter(app *models.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		filterType := r.URL.Query().Get("filter")

		var query string
		var args []interface{}

		if category != "" {
			query = `
                SELECT p.id, p.user_id, p.title, p.content, p.created_at
                FROM posts p
                JOIN post_categories pc ON p.id = pc.post_id
                JOIN categories c ON pc.category_id = c.id
                WHERE c.name = ?
                ORDER BY p.created_at DESC`
			args = append(args, category)
		} else if filterType == "created" {
			sessionID, _ := r.Cookie("session_id")
			var userID int
			err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
			if err != nil {
				utils.RenderError(app, w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			query = "SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = ? ORDER BY created_at DESC"
			args = append(args, userID)
		} else if filterType == "liked" {
			sessionID, _ := r.Cookie("session_id")
			var userID int
			err := app.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?", sessionID.Value, time.Now()).Scan(&userID)
			if err != nil {
				utils.RenderError(app, w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			query = `
                SELECT p.id, p.user_id, p.title, p.content, p.created_at
                FROM posts p
                JOIN likes l ON p.id = l.post_id
                WHERE l.user_id = ? AND l.is_like = 1
                ORDER BY p.created_at DESC`
			args = append(args, userID)
		} else {
			query = "SELECT id, user_id, title, content, created_at FROM posts ORDER BY created_at DESC"
		}

		rows, err := app.DB.Query(query, args...)
		if err != nil {
			log.Printf("Failed to query posts for filter: %v", err)
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var p models.Post
			if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
				utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			posts = append(posts, p)
		}
		log.Printf("Fetched %d posts for filter", len(posts))

		rows, err = app.DB.Query("SELECT id, name FROM categories")
		if err != nil {
			log.Printf("Failed to query categories for filter: %v", err)
			utils.RenderError(app, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var categories []struct {
			ID   int
			Name string
		}
		for rows.Next() {
			var c struct {
				ID   int
				Name string
			}
			rows.Scan(&c.ID, &c.Name)
			categories = append(categories, c)
		}

		err = app.Templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"Title":      "Home",
			"Posts":      posts,
			"Categories": categories,
			"IsLoggedIn": true,
		})
		if err != nil {
			log.Printf("Failed to render index.html for filter: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

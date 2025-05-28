package tests

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"forum/handlers"
	"forum/models"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func TestNewPostHandler(t *testing.T) {
	app := &models.App{
		DB:        db,
		Templates: template.Must(template.ParseGlob("../templates/*.html")),
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", "test", "test@example.com", string(hash))

	sessionID, _ := uuid.NewV4()
	db.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID.String(), 1, time.Now().Add(24*time.Hour))

	db.Exec("INSERT INTO categories (name) VALUES (?)", "Technology")

	req := httptest.NewRequest("POST", "/posts/new", strings.NewReader("title=Test Post&content=Test Content&categories=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID.String()})
	rr := httptest.NewRecorder()

	handlers.NewPost(app)(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM posts WHERE title = ?", "Test Post").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 1 post, got %d", count)
	}
}

func TestPostHandlerWithCommentLikes(t *testing.T) {
	app := &models.App{
		DB:        db,
		Templates: template.Must(template.ParseGlob("../templates/*.html")),
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", "test", "test@example.com", string(hash))

	db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", 1, "Test Post", "Test Content")

	db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", 1, 1, "Test Comment")

	db.Exec("INSERT INTO likes (user_id, comment_id, is_like) VALUES (?, ?, ?)", 1, 1, true)

	sessionID, _ := uuid.NewV4()
	db.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID.String(), 1, time.Now().Add(24*time.Hour))

	req := httptest.NewRequest("GET", "/posts/1", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID.String()})
	rr := httptest.NewRecorder()

	handlers.Find(app)(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Likes: 1") {
		t.Errorf("Expected comment like count in response, not found")
	}
}

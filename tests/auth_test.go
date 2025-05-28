package tests

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"forum/handlers"
	"forum/models"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	schema, _ := os.ReadFile("../database/schema.sql")
	_, err = db.Exec(string(schema))
	if err != nil {
		panic(err)
	}
}

func TestRegisterHandler(t *testing.T) {
	app := &models.App{
		DB:        db,
		Templates: template.Must(template.ParseGlob("../templates/*.html")),
	}

	req := httptest.NewRequest("POST", "/register", strings.NewReader("username=testuser&email=test@example.com&password=pass123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handlers.Register(app)(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "test@example.com").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func TestLoginHandler(t *testing.T) {
	app := &models.App{
		DB:        db,
		Templates: template.Must(template.ParseGlob("../templates/*.html")),
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", "test", "test@example.com", string(hash))

	req := httptest.NewRequest("POST", "/login", strings.NewReader("email=test@example.com&password=pass123"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handlers.Login(app)(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status %d, got %d", http.StatusSeeOther, rr.Code)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 1 session, got %d", count)
	}
}

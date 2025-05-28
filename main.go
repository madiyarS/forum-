package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"forum/handlers"
	"forum/middleware"
	"forum/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema.sql:", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to apply schema:", err)
	}

	// Verify table creation
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Fatal("Failed to query tables:", err)
	}
	defer rows.Close()
	tables := make(map[string]bool)
	tableCount := 0
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal("Failed to scan table name:", err)
		}
		tables[name] = true
		tableCount++
	}
	requiredTables := []string{"users", "posts", "comments", "categories", "post_categories", "likes", "sessions"}
	for _, table := range requiredTables {
		if !tables[table] {
			log.Fatal("Table not created:", table)
		}
	}
	log.Printf("Database schema applied successfully with %d tables", tableCount)

	// Insert categories and verify
	_, err = db.Exec(`
        INSERT OR IGNORE INTO categories (name) VALUES ('Technology');
        INSERT OR IGNORE INTO categories (name) VALUES ('Sports');
        INSERT OR IGNORE INTO categories (name) VALUES ('Entertainment');
        INSERT OR IGNORE INTO categories (name) VALUES ('General');
    `)
	if err != nil {
		log.Fatal("Failed to insert categories:", err)
	}

	var categoryCount int
	err = db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&categoryCount)
	if err != nil {
		log.Fatal("Failed to query categories:", err)
	}
	log.Printf("Categories initialized successfully: %d categories", categoryCount)

	// Parse templates with error logging
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	log.Println("Templates parsed successfully")

	app := &models.App{
		DB:        db,
		Templates: tmpl,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware.LoggingMiddleware(handlers.Index(app)))
	mux.HandleFunc("/register", middleware.LoggingMiddleware(handlers.Register(app)))
	mux.HandleFunc("/login", middleware.LoggingMiddleware(handlers.Login(app)))
	mux.HandleFunc("/logout", middleware.LoggingMiddleware(handlers.Logout(app)))
	mux.HandleFunc("/create-post", middleware.LoggingMiddleware(middleware.AuthMiddleware(app)(handlers.NewPost(app))))
	mux.HandleFunc("/post/", middleware.LoggingMiddleware(handlers.Find(app)))
	mux.HandleFunc("/post/like", middleware.LoggingMiddleware(middleware.AuthMiddleware(app)(handlers.Like(app))))
	mux.HandleFunc("/post/comment", middleware.LoggingMiddleware(middleware.AuthMiddleware(app)(handlers.Comment(app))))
	mux.HandleFunc("/posts", middleware.LoggingMiddleware(handlers.Filter(app)))
	mux.HandleFunc("/my-posts", middleware.LoggingMiddleware(handlers.Filter(app)))
	mux.HandleFunc("/liked-posts", middleware.LoggingMiddleware(handlers.Filter(app)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

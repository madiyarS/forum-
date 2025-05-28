package main

import (
	"html/template"
	"log"
	"net/http"

	"forum/internal/interfaces/database/sqlite"
	"forum/internal/interfaces/http/handlers"
	"forum/internal/interfaces/http/middleware"
	"forum/internal/usecases/auth"
	"forum/internal/usecases/comments"
	"forum/internal/usecases/likes"
	"forum/internal/usecases/posts"
)

// Main initializes and starts the server
func main() {
	// Initialize database
	db, err := sqlite.New()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Parse templates
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	log.Println("Templates parsed successfully")

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db)
	postRepo := sqlite.NewPostRepository(db)
	commentRepo := sqlite.NewCommentRepository(db)
	categoryRepo := sqlite.NewCategoryRepository(db)
	likeRepo := sqlite.NewLikeRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)

	// Initialize use cases
	authService := auth.NewService(userRepo, sessionRepo)
	postService := posts.NewService(postRepo, categoryRepo)
	commentService := comments.NewService(commentRepo)
	likeService := likes.NewService(likeRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, tmpl)
	postHandler := handlers.NewPostHandler(postService, tmpl)
	commentHandler := handlers.NewCommentHandler(commentService)
	likeHandler := handlers.NewLikeHandler(likeService)

	// Set up HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware.LoggingMiddleware(postHandler.List))
	mux.HandleFunc("/register", middleware.LoggingMiddleware(authHandler.Register))
	mux.HandleFunc("/login", middleware.LoggingMiddleware(authHandler.Login))
	mux.HandleFunc("/logout", middleware.LoggingMiddleware(authHandler.Logout))
	mux.HandleFunc("/create-post", middleware.LoggingMiddleware(middleware.AuthMiddleware(authService)(http.HandlerFunc(postHandler.Create))))
	mux.HandleFunc("/post/", middleware.LoggingMiddleware(postHandler.Get))
	mux.HandleFunc("/post/like", middleware.LoggingMiddleware(middleware.AuthMiddleware(authService)(http.HandlerFunc(likeHandler.LikePost))))
	mux.HandleFunc("/post/comment", middleware.LoggingMiddleware(middleware.AuthMiddleware(authService)(http.HandlerFunc(commentHandler.Create))))
	mux.HandleFunc("/comment/like", middleware.LoggingMiddleware(middleware.AuthMiddleware(authService)(http.HandlerFunc(likeHandler.LikeComment))))
	mux.HandleFunc("/posts", middleware.LoggingMiddleware(postHandler.Filter))
	mux.HandleFunc("/my-posts", middleware.LoggingMiddleware(middleware.AuthMiddleware(authService)(http.HandlerFunc(postHandler.Filter))))
	mux.HandleFunc("/liked-posts", middleware.LoggingMiddleware(middleware.AuthMiddleware(authService)(http.HandlerFunc(postHandler.Filter))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/usecases/posts"
)

// PostHandler handles HTTP requests for posts
type PostHandler struct {
	service *posts.Service
	tmpl    *template.Template
}

// NewPostHandler creates a new post handler
func NewPostHandler(service *posts.Service, tmpl *template.Template) *PostHandler {
	return &PostHandler{
		service: service,
		tmpl:    tmpl,
	}
}

// List handles listing all posts
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.List()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	categories, err := h.service.FindAllCategories()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Title":      "Home",
		"Posts":      posts,
		"Categories": categories,
		"IsLoggedIn": r.Context().Value("isLoggedIn") == true,
	})
}

// Create handles creating a new post
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.Context().Value("userID").(int)
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		categories := r.Form["categories"]

		var categoryIDs []int
		for _, cat := range categories {
			id, _ := strconv.Atoi(cat)
			categoryIDs = append(categoryIDs, id)
		}

		_, err := h.service.Create(userID, title, content, categoryIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	categories, err := h.service.FindAllCategories()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "create-post.html", map[string]interface{}{
		"Title":      "New Post",
		"Categories": categories,
	})
}

// Get handles retrieving a post
func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	post, err := h.service.Get(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.tmpl.ExecuteTemplate(w, "post.html", map[string]interface{}{
		"Title":        "Post",
		"Post":         post,
		"Comments":     post.Comments,
		"LikeCount":    post.LikeCount,
		"DislikeCount": post.DislikeCount,
		"IsLoggedIn":   r.Context().Value("isLoggedIn") == true,
		"HasLiked":     r.Context().Value("hasLiked") == true,
		"HasDisliked":  r.Context().Value("hasDisliked") == true,
	})
}

// Filter handles filtering posts
func (h *PostHandler) Filter(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	userID := 0
	if r.Context().Value("userID") != nil {
		userID = r.Context().Value("userID").(int)
	}
	filterType := ""
	if r.URL.Path == "/my-posts" {
		filterType = "created"
	} else if r.URL.Path == "/liked-posts" {
		filterType = "liked"
	}

	posts, err := h.service.Filter(category, userID, filterType)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	categories, err := h.service.FindAllCategories()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Title":      "Filtered Posts",
		"Posts":      posts,
		"Categories": categories,
		"IsLoggedIn": r.Context().Value("isLoggedIn") == true,
	})
}
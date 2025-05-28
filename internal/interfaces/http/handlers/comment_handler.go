package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/internal/usecases/comments"
)

// CommentHandler handles HTTP requests for comments
type CommentHandler struct {
	service *comments.Service
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(service *comments.Service) *CommentHandler {
	return &CommentHandler{
		service: service,
	}
}

// Create handles creating a new comment
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)
	content := strings.TrimSpace(r.FormValue("content"))

	_, err = h.service.Create(postID, userID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}
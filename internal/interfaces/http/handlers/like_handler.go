package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/usecases/likes"
)

// LikeHandler handles HTTP requests for likes
type LikeHandler struct {
	service *likes.Service
}

// NewLikeHandler creates a new like handler
func NewLikeHandler(service *likes.Service) *LikeHandler {
	return &LikeHandler{
		service: service,
	}
}

// LikePost handles liking a post
func (h *LikeHandler) LikePost(w http.ResponseWriter, r *http.Request) {
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
	isLike := r.FormValue("is_like") == "true"

	err = h.service.LikePost(userID, postID, isLike)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LikeComment handles liking a comment
func (h *LikeHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)
	isLike := r.FormValue("is_like") == "true"

	err = h.service.LikeComment(userID, commentID, isLike)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
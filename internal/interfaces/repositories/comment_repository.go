package repositories

import "forum/internal/entities"

// CommentRepository defines methods for comment persistence
type CommentRepository interface {
	Create(postID, userID int, content string) (*entities.Comment, error)
}
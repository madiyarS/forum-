package comments

import "forum/internal/entities"

// CommentRepository defines the interface for comment storage
type CommentRepository interface {
	Create(postID, userID int, content string) (*entities.Comment, error)
}

// Service handles comment use cases
type Service struct {
	commentRepo CommentRepository
}

// NewService creates a new comment service
func NewService(commentRepo CommentRepository) *Service {
	return &Service{
		commentRepo: commentRepo,
	}
}
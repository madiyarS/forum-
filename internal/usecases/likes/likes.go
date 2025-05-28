package likes

import "forum/internal/entities"

// LikeRepository defines the interface for like storage
type LikeRepository interface {
	CreateOrUpdate(userID int, postID, commentID *int, isLike bool) error
	FindByUserAndTarget(userID int, postID, commentID *int) (*entities.Like, error)
}

// Service handles like use cases
type Service struct {
	likeRepo LikeRepository
}

// NewService creates a new like service
func NewService(likeRepo LikeRepository) *Service {
	return &Service{
		likeRepo: likeRepo,
	}
}

// GetUserLike retrieves a user's like status for a post or comment
func (s *Service) GetUserLike(userID int, postID *int, commentID *int) (*entities.Like, error) {
	return s.likeRepo.FindByUserAndTarget(userID, postID, commentID)
}

package likes

// LikeRepository defines the interface for like storage
type LikeRepository interface {
	CreateOrUpdate(userID int, postID, commentID *int, isLike bool) error
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
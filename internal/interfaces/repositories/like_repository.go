package repositories

// LikeRepository defines methods for like persistence
type LikeRepository interface {
	CreateOrUpdate(userID int, postID, commentID *int, isLike bool) error
}
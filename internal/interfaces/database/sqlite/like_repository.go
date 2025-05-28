package sqlite

import (
	"database/sql"
)

// LikeRepository implements repositories.LikeRepository
type LikeRepository struct {
	db *sql.DB
}

// NewLikeRepository creates a new like repository
func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// CreateOrUpdate adds or updates a like
func (r *LikeRepository) CreateOrUpdate(userID int, postID, commentID *int, isLike bool) error {
	_, err := r.db.Exec("INSERT OR REPLACE INTO likes (user_id, post_id, comment_id, is_like) VALUES (?, ?, ?, ?)", userID, postID, commentID, isLike)
	return err
}
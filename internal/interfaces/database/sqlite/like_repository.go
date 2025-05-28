package sqlite

import (
	"database/sql"
	"forum/internal/entities"
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

// FindByUserAndTarget retrieves a like by user and target (post or comment)
func (r *LikeRepository) FindByUserAndTarget(userID int, postID, commentID *int) (*entities.Like, error) {
	var like entities.Like
	var query string
	var args []interface{}

	if postID != nil {
		query = "SELECT id, user_id, post_id, comment_id, is_like FROM likes WHERE user_id = ? AND post_id = ?"
		args = []interface{}{userID, *postID}
	} else if commentID != nil {
		query = "SELECT id, user_id, post_id, comment_id, is_like FROM likes WHERE user_id = ? AND comment_id = ?"
		args = []interface{}{userID, *commentID}
	} else {
		return nil, nil
	}

	err := r.db.QueryRow(query, args...).Scan(&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.IsLike)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &like, nil
}

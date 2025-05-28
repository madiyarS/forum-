package sqlite

import (
	"database/sql"
	"time"

	"forum/internal/entities"
)

// CommentRepository implements repositories.CommentRepository
type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create adds a new comment
func (r *CommentRepository) Create(postID, userID int, content string) (*entities.Comment, error) {
	result, err := r.db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entities.Comment{
		ID:        int(id),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}
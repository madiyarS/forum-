package entities

import "time"

// Like represents a like or dislike on a post or comment
type Like struct {
	ID        int
	UserID    int
	PostID    *int
	CommentID *int
	IsLike    bool
	CreatedAt time.Time
}
package entities

import "time"

// Comment represents a comment on a post
type Comment struct {
	ID           int
	PostID       int
	UserID       int
	Content      string
	CreatedAt    time.Time
	LikeCount    int
	DislikeCount int
}
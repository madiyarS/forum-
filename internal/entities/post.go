package entities

import "time"

// Post represents a forum post
type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	CreatedAt    time.Time
	Categories   []Category
	Comments     []Comment
	LikeCount    int
	DislikeCount int
}
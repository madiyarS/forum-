package entities

import "time"

// Session represents a user session
type Session struct {
	ID         string
	UserID     int
	ExpiresAt  time.Time
	CreatedAt  time.Time
}
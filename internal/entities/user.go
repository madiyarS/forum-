package entities

import "time"

// User represents a forum user
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
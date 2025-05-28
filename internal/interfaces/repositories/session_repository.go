package repositories

import (
	"forum/internal/entities"
	"time"
)

// SessionRepository defines methods for session persistence
type SessionRepository interface {
	Create(sessionID string, userID int, expiresAt time.Time) error
	FindByID(sessionID string) (*entities.Session, error)
	DeleteByUserID(userID int) error
}

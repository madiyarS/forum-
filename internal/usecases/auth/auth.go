package auth

import (
	"forum/internal/entities"
	"time"
)

// UserRepository defines the interface for user storage
type UserRepository interface {
	Create(username, email, passwordHash string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindByUsernameOrEmail(username, email string) (*entities.User, error)
}

// SessionRepository defines the interface for session storage
type SessionRepository interface {
	Create(sessionID string, userID int, expiresAt time.Time) error
	FindByID(sessionID string) (*entities.Session, error)
	DeleteByUserID(userID int) error
}

// Service handles authentication use cases
type Service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

// NewService creates a new auth service
func NewService(userRepo UserRepository, sessionRepo SessionRepository) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// GetSessionRepo returns the session repository
func (s *Service) GetSessionRepo() SessionRepository {
	return s.sessionRepo
}
package auth

import (
	"errors"
	"forum/internal/entities"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and creates a session
func (s *Service) Login(email, password string) (*entities.Session, error) {
	// Validate inputs
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Delete existing sessions
	if err := s.sessionRepo.DeleteByUserID(user.ID); err != nil {
		return nil, err
	}

	// Create new session
	sessionID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.sessionRepo.Create(sessionID.String(), user.ID, expiresAt); err != nil {
		return nil, err
	}

	return &entities.Session{
		ID:        sessionID.String(),
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}
package auth

import (
	"errors"
	"forum/internal/entities"

	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user
func (s *Service) Register(username, email, password string) (*entities.User, error) {
	// Validate inputs
	if username == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	// Check for existing user
	existingUser, err := s.userRepo.FindByUsernameOrEmail(username, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("username or email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.userRepo.Create(username, email, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}
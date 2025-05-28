package repositories

import "forum/internal/entities"

// UserRepository defines methods for user persistence
type UserRepository interface {
	Create(username, email, passwordHash string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindByUsernameOrEmail(username, email string) (*entities.User, error)
}
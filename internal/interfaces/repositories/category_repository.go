package repositories

import "forum/internal/entities"

// CategoryRepository defines methods for category persistence
type CategoryRepository interface {
	FindAll() ([]*entities.Category, error)
	FindByName(name string) (*entities.Category, error)
}
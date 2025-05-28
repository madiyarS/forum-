package posts

import "forum/internal/entities"

// PostRepository defines the interface for post storage
type PostRepository interface {
	Create(userID int, title, content string) (*entities.Post, error)
	FindAll() ([]*entities.Post, error)
	FindByID(id int) (*entities.Post, error)
	FindByCategory(categoryName string) ([]*entities.Post, error)
	FindByUserID(userID int) ([]*entities.Post, error)
	FindByUserLikes(userID int) ([]*entities.Post, error)
	AddCategory(postID, categoryID int) error
}

// CategoryRepository defines the interface for category storage
type CategoryRepository interface {
	FindAll() ([]*entities.Category, error)
	FindByName(name string) (*entities.Category, error)
}

// Service handles post use cases
type Service struct {
	postRepo     PostRepository
	categoryRepo CategoryRepository
}

// NewService creates a new post service
func NewService(postRepo PostRepository, categoryRepo CategoryRepository) *Service {
	return &Service{
		postRepo:     postRepo,
		categoryRepo: categoryRepo,
	}
}

// FindAllCategories retrieves all categories
func (s *Service) FindAllCategories() ([]*entities.Category, error) {
	return s.categoryRepo.FindAll()
}
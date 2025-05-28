package repositories

import "forum/internal/entities"

// PostRepository defines methods for post persistence
type PostRepository interface {
	Create(userID int, title, content string) (*entities.Post, error)
	FindAll() ([]*entities.Post, error)
	FindByID(id int) (*entities.Post, error)
	FindByCategory(categoryID int) ([]*entities.Post, error)
	FindByUserID(userID int) ([]*entities.Post, error)
	FindByUserLikes(userID int) ([]*entities.Post, error)
	AddCategory(postID, categoryID int) error
}
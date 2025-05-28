package posts

import (
	"errors"
	"forum/internal/entities"
)

// Filter retrieves posts by category, user, or likes
func (s *Service) Filter(categoryName string, userID int, filterType string) ([]*entities.Post, error) {
	var posts []*entities.Post
	var err error

	switch {
	case categoryName != "":
		// Find category ID by name
		category, err := s.categoryRepo.FindByName(categoryName)
		if err != nil || category == nil {
			return nil, errors.New("category not found")
		}
		posts, err = s.postRepo.FindByCategory(categoryName)
	case filterType == "created":
		posts, err = s.postRepo.FindByUserID(userID)
	case filterType == "liked":
		posts, err = s.postRepo.FindByUserLikes(userID)
	default:
		posts, err = s.postRepo.FindAll()
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}
package posts

import (
	"errors"
	"forum/internal/entities"
)

// Create makes a new post
func (s *Service) Create(userID int, title, content string, categoryIDs []int) (*entities.Post, error) {
	// Validate inputs
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}

	// Create post
	post, err := s.postRepo.Create(userID, title, content)
	if err != nil {
		return nil, err
	}

	// Add categories
	for _, catID := range categoryIDs {
		if err := s.postRepo.AddCategory(post.ID, catID); err != nil {
			return nil, err
		}
	}

	return post, nil
}
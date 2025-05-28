package posts

import (
	"errors"
	"forum/internal/entities"
)

// Get retrieves a post by ID
func (s *Service) Get(id int) (*entities.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}
package posts

import "forum/internal/entities"

// List retrieves all posts
func (s *Service) List() ([]*entities.Post, error) {
	posts, err := s.postRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return posts, nil
}
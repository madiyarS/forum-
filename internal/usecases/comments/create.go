package comments

import (
	"errors"
	"forum/internal/entities"
)

// Create adds a new comment to a post
func (s *Service) Create(postID, userID int, content string) (*entities.Comment, error) {
	if content == "" {
		return nil, errors.New("comment content is required")
	}

	comment, err := s.commentRepo.Create(postID, userID, content)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
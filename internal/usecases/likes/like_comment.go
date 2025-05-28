package likes

import "errors"

// LikeComment likes or dislikes a comment
func (s *Service) LikeComment(userID, commentID int, isLike bool) error {
	if commentID <= 0 {
		return errors.New("invalid comment ID")
	}

	return s.likeRepo.CreateOrUpdate(userID, nil, &commentID, isLike)
}
package likes

import "errors"

// LikePost likes or dislikes a post
func (s *Service) LikePost(userID, postID int, isLike bool) error {
	if postID <= 0 {
		return errors.New("invalid post ID")
	}

	return s.likeRepo.CreateOrUpdate(userID, &postID, nil, isLike)
}
package auth

import "errors"

// Logout deletes a user session
func (s *Service) Logout(sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID is required")
	}

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil || session == nil {
		return errors.New("invalid session")
	}

	return s.sessionRepo.DeleteByUserID(session.UserID)
}
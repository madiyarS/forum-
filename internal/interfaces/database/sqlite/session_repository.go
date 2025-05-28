package sqlite

import (
	"database/sql"
	"time"

	"forum/internal/entities"
)

// SessionRepository implements repositories.SessionRepository
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create adds a new session
func (r *SessionRepository) Create(sessionID string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	return err
}

// FindByID retrieves a session by ID
func (r *SessionRepository) FindByID(sessionID string) (*entities.Session, error) {
	var session entities.Session
	err := r.db.QueryRow("SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?", sessionID).
		Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteByUserID removes all sessions for a user
func (r *SessionRepository) DeleteByUserID(userID int) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
package sqlite

import (
	"database/sql"
	"time"

	"forum/internal/entities"
)

// UserRepository implements repositories.UserRepository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create adds a new user
func (r *UserRepository) Create(username, email, passwordHash string) (*entities.User, error) {
	result, err := r.db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, email, passwordHash)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:           int(id),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail retrieves a user by username or email
func (r *UserRepository) FindByUsernameOrEmail(username, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE username = ? OR email = ?", username, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
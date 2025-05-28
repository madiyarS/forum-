package tests

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDatabaseOperations(t *testing.T) {
	// Test CREATE and INSERT
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", "test", "test@example.com", string(hash))
	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}

	// Test SELECT
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "test@example.com").Scan(&count)
	if err != nil || count != 1 {
		t.Errorf("Expected 1 user, got %d, error: %v", count, err)
	}
}

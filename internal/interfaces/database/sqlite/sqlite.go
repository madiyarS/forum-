package sqlite

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// New initializes a new SQLite database
func New() (*sql.DB, error) {
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return nil, err
	}

	schema, err := os.ReadFile("internal/interfaces/database/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO categories (name) VALUES ('Technology');
		INSERT OR IGNORE INTO categories (name) VALUES ('Sports');
		INSERT OR IGNORE INTO categories (name) VALUES ('Entertainment');
		INSERT OR IGNORE INTO categories (name) VALUES ('General');
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
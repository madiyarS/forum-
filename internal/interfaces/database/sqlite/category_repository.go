package sqlite

import (
	"database/sql"
	"forum/internal/entities"
)

// CategoryRepository implements repositories.CategoryRepository
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// FindAll retrieves all categories
func (r *CategoryRepository) FindAll() ([]*entities.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.Category
	for rows.Next() {
		var c entities.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	return categories, nil
}

// FindByName retrieves a category by name
func (r *CategoryRepository) FindByName(name string) (*entities.Category, error) {
	var c entities.Category
	err := r.db.QueryRow("SELECT id, name FROM categories WHERE name = ?", name).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
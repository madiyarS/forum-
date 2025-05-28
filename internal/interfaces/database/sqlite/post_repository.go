package sqlite

import (
	"database/sql"
	"time"

	"forum/internal/entities"
)

// PostRepository implements repositories.PostRepository
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create adds a new post
func (r *PostRepository) Create(userID int, title, content string) (*entities.Post, error) {
	result, err := r.db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entities.Post{
		ID:        int(id),
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

// FindAll retrieves all posts
func (r *PostRepository) FindAll() ([]*entities.Post, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		var p entities.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

// FindByID retrieves a post by ID
func (r *PostRepository) FindByID(id int) (*entities.Post, error) {
	var p entities.Post
	err := r.db.QueryRow("SELECT id, user_id, title, content, created_at FROM posts WHERE id = ?", id).
		Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch categories
	rows, err := r.db.Query("SELECT c.id, c.name FROM categories c JOIN post_categories pc ON c.id = pc.category_id WHERE pc.post_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var c entities.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	p.Categories = categories

	// Fetch comments
	commentRows, err := r.db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at ASC", id)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	var comments []entities.Comment
	for commentRows.Next() {
		var c entities.Comment
		if err := commentRows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		// Fetch like/dislike counts for comment
		err = r.db.QueryRow("SELECT COUNT(CASE WHEN is_like = 1 THEN 1 END), COUNT(CASE WHEN is_like = 0 THEN 1 END) FROM likes WHERE comment_id = ?", c.ID).
			Scan(&c.LikeCount, &c.DislikeCount)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	p.Comments = comments

	// Fetch like/dislike counts for post
	err = r.db.QueryRow("SELECT COUNT(CASE WHEN is_like = 1 THEN 1 END), COUNT(CASE WHEN is_like = 0 THEN 1 END) FROM likes WHERE post_id = ?", id).
		Scan(&p.LikeCount, &p.DislikeCount)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// FindByCategory retrieves posts by category name
func (r *PostRepository) FindByCategory(categoryName string) ([]*entities.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at
		FROM posts p
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name = ?
		ORDER BY p.created_at DESC`, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		var p entities.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

// FindByUserID retrieves posts by user
func (r *PostRepository) FindByUserID(userID int) ([]*entities.Post, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		var p entities.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

// FindByUserLikes retrieves posts liked by a user
func (r *PostRepository) FindByUserLikes(userID int) ([]*entities.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at
		FROM posts p
		JOIN likes l ON p.id = l.post_id
		WHERE l.user_id = ? AND l.is_like = 1
		ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		var p entities.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

// AddCategory associates a category with a post
func (r *PostRepository) AddCategory(postID, categoryID int) error {
	_, err := r.db.Exec("INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	return err
}
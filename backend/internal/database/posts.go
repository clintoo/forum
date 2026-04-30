package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/backend/internal/models"
)

// GetAllPosts returns all posts ordered by creation date (newest first)
func (db *SQLiteStore) GetAllPosts() ([]models.Post, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at
        FROM posts p
        ORDER BY p.created_at DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	var postIDs []int
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	if len(postIDs) > 0 {
		// Bulk fetch categories
		categoriesMap, err := db.getPostsCategoriesBulk(postIDs)
		if err == nil {
			for i := range posts {
				posts[i].Categories = categoriesMap[posts[i].ID]
			}
		}

		// Bulk fetch stats
		statsMap, err := db.GetPostsStats(postIDs)
		if err == nil {
			for i := range posts {
				if stats, ok := statsMap[posts[i].ID]; ok {
					posts[i].Likes = stats.Likes
					posts[i].Dislikes = stats.Dislikes
					posts[i].CommentCount = stats.CommentCount
				}
			}
		}
	}

	return posts, nil
}

// getPostsCategoriesBulk fetches categories for multiple posts efficiently
func (db *SQLiteStore) getPostsCategoriesBulk(postIDs []int) (map[int][]string, error) {
	if len(postIDs) == 0 {
		return make(map[int][]string), nil
	}

	placeholders := make([]string, len(postIDs))
	args := make([]any, len(postIDs))
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT pc.post_id, c.name
        FROM categories c
        JOIN post_categories pc ON c.id = pc.category_id
        WHERE pc.post_id IN (%s)
    `, strings.Join(placeholders, ", "))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]string)
	for rows.Next() {
		var postID int
		var name string
		if err := rows.Scan(&postID, &name); err != nil {
			return nil, err
		}
		result[postID] = append(result[postID], name)
	}
	return result, nil
}

// GetPostCategories returns category names for a post
func (db *SQLiteStore) GetPostCategories(postID int) ([]string, error) {
	query := `
        SELECT c.name
        FROM categories c
        JOIN post_categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
    `

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		categories = append(categories, name)
	}
	return categories, nil
}

// CreatePost inserts a new post with categories and returns the full post model
func (db *SQLiteStore) CreatePost(post *models.Post) (*models.Post, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the post
	query := `INSERT INTO posts (user_id, title, content, image_path) VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(query, post.UserID, post.Title, post.Content, post.ImagePath)
	if err != nil {
		return nil, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Insert categories
	for _, categoryID := range post.CategoryIDs {
		_, err = tx.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
		if err != nil {
			return nil, err
		}
	}

	// Fetch category names
	var categoryNames []string
	if len(post.CategoryIDs) > 0 {
		rows, err := tx.Query(`SELECT name FROM categories WHERE id IN (SELECT category_id FROM post_categories WHERE post_id = ?)`, postID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err == nil {
					categoryNames = append(categoryNames, name)
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	now := time.Now()
	return &models.Post{
		ID:         int(postID),
		UserID:     post.UserID,
		Title:      post.Title,
		Content:    post.Content,
		ImagePath:  post.ImagePath,
		CreatedAt:  now,
		UpdatedAt:  now,
		Categories: categoryNames,
	}, nil
}

// GetPostByID retrieves a specific post by ID
func (db *SQLiteStore) GetPostByID(postID int) (*models.Post, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at
        FROM posts p
        WHERE p.id = ?
    `

	var post models.Post
	err := db.QueryRow(query, postID).Scan(&post.ID, &post.UserID, &post.Title,
		&post.Content, &post.ImagePath, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	categories, _ := db.GetPostCategories(postID)
	post.Categories = categories

	// Fetch single post stats
	statsMap, err := db.GetPostsStats([]int{postID})
	if err == nil {
		if stats, ok := statsMap[postID]; ok {
			post.Likes = stats.Likes
			post.Dislikes = stats.Dislikes
			post.CommentCount = stats.CommentCount
		}
	}

	return &post, nil
}

// UpdatePost updates a post's title, content, image path, and categories atomically.
func (db *SQLiteStore) UpdatePost(post *models.Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`UPDATE posts SET title = ?, content = ?, image_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		post.Title, post.Content, post.ImagePath, post.ID,
	)
	if err != nil {
		return err
	}

	if post.CategoryIDs != nil {
		if _, err = tx.Exec(`DELETE FROM post_categories WHERE post_id = ?`, post.ID); err != nil {
			return err
		}
		for _, catID := range post.CategoryIDs {
			if _, err = tx.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, post.ID, catID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// DeletePost deletes a post (cascade will delete comments and votes)
func (db *SQLiteStore) DeletePost(postID int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := db.Exec(query, postID)
	return err
}

// GetPostStats returns like, dislike, and comment counts for a post.
// post_votes is scanned once using conditional aggregation.
func (db *SQLiteStore) GetPostStats(postID int) (likes, dislikes, comments int, err error) {
	query := `
        SELECT
            COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
            COALESCE(SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes,
            (SELECT COUNT(*) FROM comments WHERE post_id = ?) as comments
        FROM post_votes
        WHERE post_id = ?
    `
	err = db.QueryRow(query, postID, postID).Scan(&likes, &dislikes, &comments)
	return
}

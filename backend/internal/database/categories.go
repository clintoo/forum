package database

import (
	"forum/backend/internal/models"
)

// GetAllCategories returns all categories
func (db *SQLiteStore) GetAllCategories() ([]models.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetPostsByCategory returns all posts in a specific category
func (db *SQLiteStore) GetPostsByCategory(categoryID int) ([]models.Post, error) {
	query := `
        SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at
        FROM posts p
        JOIN post_categories pc ON p.id = pc.post_id
        WHERE pc.category_id = ?
        ORDER BY p.created_at DESC
    `

	rows, err := db.Query(query, categoryID)
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

// CreateCategory creates a new category (admin function) and returns the full model
func (db *SQLiteStore) CreateCategory(name string) (*models.Category, error) {
	query := `INSERT INTO categories (name) VALUES (?)`
	result, err := db.Exec(query, name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &models.Category{
		ID:   int(id),
		Name: name,
	}, nil
}

// GetCategoryByID retrieves a specific category by ID
func (db *SQLiteStore) GetCategoryByID(id int) (*models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = ?`
	var category models.Category
	err := db.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

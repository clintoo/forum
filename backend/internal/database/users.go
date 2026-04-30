package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/backend/internal/models"
)

// GetAllUsers returns all users
func (db *SQLiteStore) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, email, username, avatar_path, created_at FROM users ORDER BY created_at DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.AvatarPath, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// CreateUser inserts a new user and returns the full user model
func (db *SQLiteStore) CreateUser(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (email, username, avatar_path, password) VALUES (?, ?, ?, ?)`

	result, err := db.DB.Exec(query, user.Email, user.Username, user.AvatarPath, user.Password)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:         int(userID),
		Email:      user.Email,
		Username:   user.Username,
		AvatarPath: user.AvatarPath,
		CreatedAt:  time.Now(),
	}, nil
}

// GetUserByID retrieves a user by ID
func (db *SQLiteStore) GetUserByID(userID int) (*models.User, error) {
	query := `SELECT id, email, username, avatar_path, created_at FROM users WHERE id = ?`

	var user models.User
	err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.Email, &user.Username, &user.AvatarPath, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email (for login)
func (db *SQLiteStore) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, username, avatar_path, password, created_at FROM users WHERE email = ?`

	var user models.User
	err := db.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Username, &user.AvatarPath, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username (for registration check)
func (db *SQLiteStore) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, email, username, avatar_path, created_at FROM users WHERE username = ?`

	var user models.User
	err := db.DB.QueryRow(query, username).Scan(&user.ID, &user.Email, &user.Username, &user.AvatarPath, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetPostsByUserID gets all posts by a user
func (db *SQLiteStore) GetPostsByUserID(userID int) ([]models.Post, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at
        FROM posts p
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC
    `

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// DeleteUser removes a user and all their associated data by ID.
func (db *SQLiteStore) DeleteUser(userID int) error {
	_, err := db.DB.Exec(`DELETE FROM users WHERE id = ?`, userID)
	return err
}

// UpdateUser updates the mutable fields of a user. Only non-nil pointers are applied.
func (db *SQLiteStore) UpdateUser(userID int, username, email, avatarPath, passwordHash *string) error {
	query := `UPDATE users SET
		username    = COALESCE(?, username),
		email       = COALESCE(?, email),
		avatar_path = COALESCE(?, avatar_path),
		password    = COALESCE(?, password)
	WHERE id = ?`
	_, err := db.DB.Exec(query, username, email, avatarPath, passwordHash, userID)
	return err
}

func (db *SQLiteStore) GetUserComments(userID int) ([]models.Comment, error) {
	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at
        FROM comments c
        WHERE c.user_id = ?
        ORDER BY c.created_at DESC
    `

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetUsersByIDs returns a map of userID -> User for the given IDs in a single query.
// Duplicate IDs are deduplicated before querying.
func (db *SQLiteStore) GetUsersByIDs(userIDs []int) (map[int]models.User, error) {
	if len(userIDs) == 0 {
		return map[int]models.User{}, nil
	}

	seen := make(map[int]struct{}, len(userIDs))
	unique := make([]int, 0, len(userIDs))
	for _, id := range userIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			unique = append(unique, id)
		}
	}

	placeholders := make([]string, len(unique))
	args := make([]any, len(unique))
	for i, id := range unique {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT id, email, username, avatar_path, created_at FROM users WHERE id IN (%s)`,
		strings.Join(placeholders, ", "),
	)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]models.User, len(unique))
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.AvatarPath, &u.CreatedAt); err != nil {
			return nil, err
		}
		result[u.ID] = u
	}
	return result, nil
}

// GetUserLikes gets all posts liked by a user
func (db *SQLiteStore) GetUserLikes(userID int) ([]models.Post, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at
        FROM post_votes pv
        JOIN posts p ON pv.post_id = p.id
        WHERE pv.user_id = ? AND pv.vote_type = 1
        ORDER BY pv.created_at DESC
    `

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ImagePath,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

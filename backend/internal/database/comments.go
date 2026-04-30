package database

import (
	"database/sql"
	"time"

	"forum/backend/internal/models"
)

// GetCommentsByPostID returns all comments on a specific post
func (db *SQLiteStore) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at
        FROM comments c
        WHERE c.post_id = ?
        ORDER BY c.created_at ASC
    `

	rows, err := db.Query(query, postID)
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

// CreateComment inserts a new comment on a post and returns the full comment model
func (db *SQLiteStore) CreateComment(comment *models.Comment) (*models.Comment, error) {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	result, err := db.Exec(query, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return nil, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &models.Comment{
		ID:        int(commentID),
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateComment updates a specific comment
func (db *SQLiteStore) UpdateComment(comment *models.Comment) error {
	query := `UPDATE comments SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, comment.Content, comment.ID)
	return err
}

// DeleteComment deletes a specific comment
func (db *SQLiteStore) DeleteComment(commentID int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := db.Exec(query, commentID)
	return err
}

// GetCommentByID retrieves a specific comment with full details
func (db *SQLiteStore) GetCommentByID(commentID int) (*models.Comment, error) {
	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at
        FROM comments c
        WHERE c.id = ?
    `

	var comment models.Comment
	err := db.QueryRow(query, commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID,
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &comment, nil
}

// GetCommentStats returns like and dislike counts for a comment
func (db *SQLiteStore) GetCommentStats(commentID int) (likes, dislikes int, err error) {
	query := `
        SELECT
            COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
            COALESCE(SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes
        FROM comment_votes
        WHERE comment_id = ?
    `
	err = db.DB.QueryRow(query, commentID).Scan(&likes, &dislikes)
	return
}

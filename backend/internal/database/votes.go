package database

import (
	"database/sql"
	"fmt"
	"strings"

	"forum/backend/internal/models"
)

// VoteOnPost adds, updates, or removes a vote on a post
// voteType: 1 = like, -1 = dislike, 0 = remove vote
func (db *SQLiteStore) VoteOnPost(userID, postID, voteType int) error {
	if voteType == 0 {
		// Remove vote
		query := `DELETE FROM post_votes WHERE user_id = ? AND post_id = ?`
		_, err := db.DB.Exec(query, userID, postID)
		return err
	}

	// Insert or update vote
	query := `
        INSERT INTO post_votes (user_id, post_id, vote_type)
        VALUES (?, ?, ?)
        ON CONFLICT(user_id, post_id) DO UPDATE SET vote_type = ?
    `
	_, err := db.DB.Exec(query, userID, postID, voteType, voteType)
	return err
}

// VoteOnComment adds, updates, or removes a vote on a comment
func (db *SQLiteStore) VoteOnComment(userID, commentID, voteType int) error {
	if voteType == 0 {
		query := `DELETE FROM comment_votes WHERE user_id = ? AND comment_id = ?`
		_, err := db.DB.Exec(query, userID, commentID)
		return err
	}

	query := `
        INSERT INTO comment_votes (user_id, comment_id, vote_type)
        VALUES (?, ?, ?)
        ON CONFLICT(user_id, comment_id) DO UPDATE SET vote_type = ?
    `
	_, err := db.DB.Exec(query, userID, commentID, voteType, voteType)
	return err
}

// GetPostVoteCount returns like and dislike counts for a post
func (db *SQLiteStore) GetPostVoteCount(postID int) (likes, dislikes int, err error) {
	query := `
        SELECT
            COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
            COALESCE(SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes
        FROM post_votes
        WHERE post_id = ?
    `
	err = db.DB.QueryRow(query, postID).Scan(&likes, &dislikes)
	return
}

// GetUserPostVote returns the vote type for a user on a post
func (db *SQLiteStore) GetUserPostVote(userID, postID int) (voteType int, err error) {
	query := `SELECT vote_type FROM post_votes WHERE user_id = ? AND post_id = ?`
	err = db.DB.QueryRow(query, userID, postID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil // No vote found
	}
	return
}

// GetUserCommentVote returns the vote type for a user on a comment
func (db *SQLiteStore) GetUserCommentVote(userID, commentID int) (voteType int, err error) {
	query := `SELECT vote_type FROM comment_votes WHERE user_id = ? AND comment_id = ?`
	err = db.DB.QueryRow(query, userID, commentID).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return
}

// GetPostsStats returns likes, dislikes, and comment counts for multiple posts in one query.
func (db *SQLiteStore) GetPostsStats(postIDs []int) (map[int]models.PostStats, error) {
	if len(postIDs) == 0 {
		return map[int]models.PostStats{}, nil
	}

	placeholders := make([]string, len(postIDs))
	args := make([]any, len(postIDs))
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT p.id,
            COALESCE(SUM(CASE WHEN pv.vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
            COALESCE(SUM(CASE WHEN pv.vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes,
            COUNT(DISTINCT c.id) as comment_count
        FROM posts p
        LEFT JOIN post_votes pv ON p.id = pv.post_id
        LEFT JOIN comments c ON p.id = c.post_id
        WHERE p.id IN (%s)
        GROUP BY p.id
    `, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]models.PostStats, len(postIDs))
	for rows.Next() {
		var postID int
		var stats models.PostStats
		if err := rows.Scan(&postID, &stats.Likes, &stats.Dislikes, &stats.CommentCount); err != nil {
			return nil, err
		}
		result[postID] = stats
	}
	return result, nil
}

// GetCommentsStats returns likes and dislikes for multiple comments in one query.
func (db *SQLiteStore) GetCommentsStats(commentIDs []int) (map[int]models.CommentStats, error) {
	if len(commentIDs) == 0 {
		return map[int]models.CommentStats{}, nil
	}

	placeholders := make([]string, len(commentIDs))
	args := make([]any, len(commentIDs))
	for i, id := range commentIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT comment_id,
            COALESCE(SUM(CASE WHEN vote_type = 1 THEN 1 ELSE 0 END), 0) as likes,
            COALESCE(SUM(CASE WHEN vote_type = -1 THEN 1 ELSE 0 END), 0) as dislikes
        FROM comment_votes
        WHERE comment_id IN (%s)
        GROUP BY comment_id
    `, strings.Join(placeholders, ", "))

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]models.CommentStats, len(commentIDs))
	for rows.Next() {
		var commentID int
		var stats models.CommentStats
		if err := rows.Scan(&commentID, &stats.Likes, &stats.Dislikes); err != nil {
			return nil, err
		}
		result[commentID] = stats
	}
	return result, nil
}

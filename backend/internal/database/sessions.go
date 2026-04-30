package database

import (
	"database/sql"
	"time"

	"forum/backend/internal/models"
	"forum/backend/internal/utils"
)

// CreateSession inserts a new session and returns the UUID
func (db *SQLiteStore) CreateSession(userID int, duration time.Duration) (*models.Session, error) {
	// Generate UUID for session ID
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(duration)

	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err = db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}, nil
}

// GetSessionByID retrieves a session by its ID
func (db *SQLiteStore) GetSessionByID(sessionID string) (*models.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = ?`

	var session models.Session
	err := db.QueryRow(query, sessionID).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// GetUserSessions gets all sessions for a user
func (db *SQLiteStore) GetUserSessions(userID int) ([]models.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE user_id = ? ORDER BY expires_at DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(&session.ID, &session.UserID, &session.ExpiresAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// DeleteSession deletes a session (logout)
func (db *SQLiteStore) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(query, sessionID)
	return err
}

// DeleteExpiredSessions cleans up expired sessions
func (db *SQLiteStore) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < ?`
	_, err := db.Exec(query, time.Now())
	return err
}

// ValidateSession checks if a session is valid and returns the user ID
func (db *SQLiteStore) ValidateSession(sessionID string) (int, bool) {
	session, err := db.GetSessionByID(sessionID)
	if err != nil || session == nil {
		return 0, false
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		db.DeleteSession(sessionID) // Clean up expired session
		return 0, false
	}

	return session.UserID, true
}

package services

import (
	"time"

	"forum/backend/internal/models"
)

type SessionStorer interface {
	CreateSession(userID int, duration time.Duration) (*models.Session, error)
	DeleteSession(sessionID string) error
	DeleteExpiredSessions() error

	GetSessionByID(sessionID string) (*models.Session, error)
	GetUserSessions(userID int) ([]models.Session, error)
	ValidateSession(sessionID string) (int, bool)
}

type SessionService struct {
	Db              SessionStorer
	SessionDuration time.Duration
}

func (s *SessionService) CreateUserSession(userID int) (*models.Session, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user id")
	}

	duration := s.SessionDuration
	if duration <= 0 {
		duration = 24 * time.Hour
	}

	session, err := s.Db.CreateSession(userID, duration)
	if err != nil {
		return nil, NewInternalError("failed to create user session", err)
	}

	return session, nil
}

func (s *SessionService) GetSessionByID(sessionID string) (*models.Session, error) {
	if sessionID == "" {
		return nil, NewValidationError("invalid session id")
	}

	session, err := s.Db.GetSessionByID(sessionID)
	if err != nil {
		return nil, NewInternalError("failed to retrieve session", err)
	}
	return session, nil
}

func (s *SessionService) GetUserSessions(userID int) ([]models.Session, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user id")
	}

	sessions, err := s.Db.GetUserSessions(userID)
	if err != nil {
		return nil, NewInternalError("failed to retrieve user sessions", err)
	}
	return sessions, nil
}

func (s *SessionService) ValidateSession(sessionID string) (int, bool) {
	if sessionID == "" {
		return 0, false
	}
	return s.Db.ValidateSession(sessionID)
}

func (s *SessionService) Logout(sessionID string) error {
	if sessionID == "" {
		return NewValidationError("invalid session token")
	}

	if err := s.Db.DeleteSession(sessionID); err != nil {
		return NewInternalError("failed to delete session", err)
	}
	return nil
}

func (s *SessionService) CleanupExpiredSessions() error {
	if err := s.Db.DeleteExpiredSessions(); err != nil {
		return NewInternalError("failed to cleanup expired sessions", err)
	}
	return nil
}

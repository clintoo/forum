package middleware

import (
	"context"
	"net/http"
)

type MiddlewareDecorator = func(http.HandlerFunc) http.HandlerFunc

type contextKey string

// UserIDKey is the context key under which the authenticated user's ID is stored.
const UserIDKey contextKey = "userID"

// GetUserID retrieves the authenticated user's ID from the request context.
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// SessionValidator is satisfied by any service that can validate a session token.
type SessionValidator interface {
	ValidateSession(sessionID string) (int, bool)
}

// AuthMiddlewareFactory returns a middleware that validates the session_token cookie and
// injects the authenticated user's ID into the request context.
func AuthMiddlewareFactory(sv SessionValidator) MiddlewareDecorator {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, ok := sv.ValidateSession(cookie.Value)
			if !ok {
				// Clear the stale cookie on the client
				http.SetCookie(w, &http.Cookie{
					Name:   "session_token",
					Value:  "",
					MaxAge: -1,
				})
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next(w, r.WithContext(ctx))
		}
	}
}

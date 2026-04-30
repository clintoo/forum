package handlers

import (
	"encoding/json"
	"net/http"

	"forum/backend/internal/dto"
	"forum/backend/internal/services"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	UserService    UserServicer
	SessionService SessionServicer
}

func (h *AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	var newUser dto.PrivateUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		return services.NewValidationError("Invalid JSON")
	}

	user, err := h.UserService.SignUp(&newUser)
	if err != nil {
		return err
	}

	session, err := h.SessionService.CreateUserSession(user.ID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})

	writeSuccess(w, http.StatusCreated, "signup successful", dto.ToPrivateUser(user))
	return nil
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	// Define struct for login credentials
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return services.NewValidationError("Invalid JSON")
	}
	// Call Userservice.Login
	user, err := h.UserService.Login(credentials.Email, credentials.Password)
	if err != nil {
		return err
	}
	// Create session
	session, err := h.SessionService.CreateUserSession(user.ID)
	if err != nil {
		return err
	}
	// Set session cookie (same as signup)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})

	writeSuccess(w, http.StatusOK, "login successful", dto.ToPrivateUser(user))
	return nil
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return services.NewValidationError("missing session token")
	}

	if err := h.SessionService.Logout(cookie.Value); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	writeSuccess(w, http.StatusOK, "logout successful", nil)
	return nil
}

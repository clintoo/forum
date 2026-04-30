package handlers

import (
	"encoding/json"
	"net/http"

	"forum/backend/internal/dto"

	"forum/backend/internal/middleware"
	"forum/backend/internal/services"
)

// UserHandler handles user profile endpoints.
type UserHandler struct {
	UserService UserServicer
}

func (h *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "users retrieved successfully", dto.ToPublicUsers(users))
	return nil
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	userID, err := getURLParamInt(r, "user_id")
	if err != nil {
		return err
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "user retrieved successfully", dto.ToPublicUser(user))
	return nil
}

func (h *UserHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	user, err := h.UserService.GetMe(userID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "profile retrieved successfully", dto.ToPrivateUser(user))
	return nil
}

func (h *UserHandler) HandleUpdateMe(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	var req struct {
		Username   *string `json:"username"`
		Email      *string `json:"email"`
		AvatarPath *string `json:"avatar_path"`
		Password   *string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	user, err := h.UserService.UpdateMe(userID, req.Username, req.Email, req.AvatarPath, req.Password)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "profile updated successfully", dto.ToPrivateUser(user))
	return nil
}

func (h *UserHandler) HandleDeleteMe(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	if err := h.UserService.DeleteMe(userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *UserHandler) HandleGetMePosts(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	posts, err := h.UserService.GetMePosts(userID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "posts retrieved successfully", dto.ToPosts(posts))
	return nil
}

func (h *UserHandler) HandleGetMeComments(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	comments, err := h.UserService.GetMeComments(userID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "comments retrieved successfully", dto.ToComments(comments))
	return nil
}

func (h *UserHandler) HandleGetMeLikes(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	likes, err := h.UserService.GetMeLikes(userID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "liked posts retrieved successfully", dto.ToPosts(likes))
	return nil
}

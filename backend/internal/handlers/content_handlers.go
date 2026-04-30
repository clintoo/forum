package handlers

import (
	"encoding/json"
	"net/http"

	"forum/backend/internal/dto"

	"forum/backend/internal/middleware"
	"forum/backend/internal/models"
	"forum/backend/internal/services"
)

type ContentHandler struct {
	PostService     PostServicer
	CommentService  CommentServicer
	ReactionService ReactionServicer
}

// --- Posts ---

func (h *ContentHandler) HandleGetPosts(w http.ResponseWriter, r *http.Request) error {
	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "posts retrieved successfully", dto.ToPosts(posts))
	return nil
}

func (h *ContentHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	var req struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		ImagePath   string `json:"image_path"`
		CategoryIDs []int  `json:"category_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	post, err := h.PostService.CreatePost(&models.Post{
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		ImagePath:   req.ImagePath,
		CategoryIDs: req.CategoryIDs,
	})
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusCreated, "post created successfully", dto.ToPost(post))
	return nil
}

func (h *ContentHandler) HandleGetPost(w http.ResponseWriter, r *http.Request) error {
	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	post, err := h.PostService.GetPostByID(postID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "post retrieved successfully", dto.ToPost(post))
	return nil
}

func (h *ContentHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	var req struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		ImagePath   string `json:"image_path"`
		CategoryIDs []int  `json:"category_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	post, err := h.PostService.UpdatePost(&models.Post{
		ID:          postID,
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		ImagePath:   req.ImagePath,
		CategoryIDs: req.CategoryIDs,
	})
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "post updated successfully", dto.ToPost(post))
	return nil
}

func (h *ContentHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	if err := h.PostService.DeletePost(postID, userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// --- Comments ---

func (h *ContentHandler) HandleGetPostComments(w http.ResponseWriter, r *http.Request) error {
	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	comments, err := h.CommentService.GetCommentsByPostID(postID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "comments retrieved successfully", dto.ToComments(comments))
	return nil
}

func (h *ContentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	comment, err := h.CommentService.CreateComment(&models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	})
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusCreated, "comment created successfully", dto.ToComment(comment))
	return nil
}

func (h *ContentHandler) HandleUpdateComment(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	commentID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	comment, err := h.CommentService.UpdateComment(&models.Comment{
		ID:      commentID,
		UserID:  userID,
		Content: req.Content,
	})
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "comment updated successfully", dto.ToComment(comment))
	return nil
}

func (h *ContentHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	commentID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	if err := h.CommentService.DeleteComment(commentID, userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// --- Reactions ---

func (h *ContentHandler) HandleReactToPost(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	var req struct {
		ReactionType string `json:"reaction_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	post, err := h.ReactionService.ReactToPost(userID, postID, req.ReactionType)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "post reaction recorded successfully", dto.ToPost(post))
	return nil
}

func (h *ContentHandler) HandleRemovePostReaction(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	post, err := h.ReactionService.RemovePostReaction(userID, postID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "post reaction removed successfully", dto.ToPost(post))
	return nil
}

func (h *ContentHandler) HandleReactToComment(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	postID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	commentID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	var req struct {
		ReactionType string `json:"reaction_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return services.NewValidationError("invalid request body")
	}

	comment, err := h.ReactionService.ReactToComment(userID, postID, commentID, req.ReactionType)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "comment reaction recorded successfully", dto.ToComment(comment))
	return nil
}

func (h *ContentHandler) HandleRemoveCommentReaction(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok || userID <= 0 {
		return services.NewValidationError("invalid user ID in context")
	}

	commentID, err := getURLParamInt(r, "id")
	if err != nil {
		return err
	}

	comment, err := h.ReactionService.RemoveCommentReaction(userID, commentID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "comment reaction removed successfully", dto.ToComment(comment))
	return nil
}

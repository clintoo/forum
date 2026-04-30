package handlers

import (
	"net/http"

	"forum/backend/internal/dto"
)

// CategoryHandler handles category endpoints.
type CategoryHandler struct {
	PostService PostServicer
}

func (h *CategoryHandler) HandleGetCategories(w http.ResponseWriter, r *http.Request) error {
	categories, err := h.PostService.GetAllCategories()
	if err != nil {
		return err
	}
	writeSuccess(w, http.StatusOK, "categories retrieved successfully", dto.ToCategories(categories))
	return nil
}

func (h *CategoryHandler) HandleGetCategoryPosts(w http.ResponseWriter, r *http.Request) error {
	categoryID, err := getURLParamInt(r, "category_id")
	if err != nil {
		return err
	}

	posts, err := h.PostService.GetPostsByCategory(categoryID)
	if err != nil {
		return err
	}

	writeSuccess(w, http.StatusOK, "posts by category retrieved successfully", dto.ToPosts(posts))
	return nil
}

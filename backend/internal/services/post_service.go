package services

import (
	"forum/backend/internal/models"
)

type PostStorer interface {
	CreatePost(post *models.Post) (*models.Post, error)
	UpdatePost(post *models.Post) error
	DeletePost(postID int) error

	GetPostByID(postID int) (*models.Post, error)
	GetAllPosts() ([]models.Post, error)
	GetPostsByUserID(userID int) ([]models.Post, error)

	GetAllCategories() ([]models.Category, error)
	GetPostsByCategory(categoryID int) ([]models.Post, error)
}

type PostService struct {
	Db PostStorer
}

func (s *PostService) CreatePost(post *models.Post) (*models.Post, error) {
	if err := post.Validate(); err != nil {
		return nil, NewValidationError(err.Error())
	}

	created, err := s.Db.CreatePost(post)
	if err != nil {
		return nil, NewInternalError("error while creating post", err)
	}
	return created, nil
}

func (s *PostService) GetPostByID(postID int) (*models.Post, error) {
	if postID <= 0 {
		return nil, NewValidationError("invalid post ID")
	}

	post, err := s.Db.GetPostByID(postID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if post == nil {
		return nil, NewNotFoundError("post not found")
	}
	return post, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	posts, err := s.Db.GetAllPosts()
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return posts, nil
}

func (s *PostService) UpdatePost(post *models.Post) (*models.Post, error) {
	if post.ID <= 0 {
		return nil, NewValidationError("invalid post ID")
	}
	if err := post.Validate(); err != nil {
		return nil, NewValidationError(err.Error())
	}

	existing, err := s.Db.GetPostByID(post.ID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if existing == nil {
		return nil, NewNotFoundError("post not found")
	}
	if existing.UserID != post.UserID {
		return nil, NewValidationError("not allowed to update post")
	}

	if err := s.Db.UpdatePost(post); err != nil {
		return nil, NewInternalError("error while updating post", err)
	}

	return s.GetPostByID(post.ID)
}

func (s *PostService) DeletePost(postID, userID int) error {
	if postID <= 0 {
		return NewValidationError("invalid post ID")
	}
	if userID <= 0 {
		return NewValidationError("invalid user ID")
	}

	post, err := s.Db.GetPostByID(postID)
	if err != nil {
		return NewInternalError("error while querying DB", err)
	}
	if post == nil {
		return NewNotFoundError("post not found")
	}
	if post.UserID != userID {
		return NewValidationError("not allowed to delete post")
	}

	if err := s.Db.DeletePost(postID); err != nil {
		return NewInternalError("error while deleting post", err)
	}
	return nil
}

func (s *PostService) GetPostsByUserID(userID int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}

	posts, err := s.Db.GetPostsByUserID(userID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return posts, nil
}

func (s *PostService) GetAllCategories() ([]models.Category, error) {
	categories, err := s.Db.GetAllCategories()
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return categories, nil
}

func (s *PostService) GetPostsByCategory(categoryID int) ([]models.Post, error) {
	if categoryID <= 0 {
		return nil, NewValidationError("invalid category ID")
	}

	posts, err := s.Db.GetPostsByCategory(categoryID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return posts, nil
}

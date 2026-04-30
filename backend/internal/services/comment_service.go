package services

import "forum/backend/internal/models"

type CommentStorer interface {
	CreateComment(comment *models.Comment) (*models.Comment, error)
	UpdateComment(comment *models.Comment) error
	DeleteComment(commentID int) error

	GetCommentByID(commentID int) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]models.Comment, error)
}

type CommentService struct {
	Db CommentStorer
}

func (s *CommentService) CreateComment(comment *models.Comment) (*models.Comment, error) {
	if err := comment.Validate(); err != nil {
		return nil, NewValidationError(err.Error())
	}
	created, err := s.Db.CreateComment(comment)
	if err != nil {
		return nil, NewInternalError("error while creating comment", err)
	}
	return created, nil
}

func (s *CommentService) GetCommentByID(commentID int) (*models.Comment, error) {
	if commentID <= 0 {
		return nil, NewValidationError("Invalid comment ID")
	}
	comment, err := s.Db.GetCommentByID(commentID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if comment == nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID int) ([]models.Comment, error) {
	if postID <= 0 {
		return nil, NewValidationError("invalid post ID")
	}

	comments, err := s.Db.GetCommentsByPostID(postID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return comments, nil
}

func (s *CommentService) UpdateComment(comment *models.Comment) (*models.Comment, error) {
	if comment.ID <= 0 {
		return nil, NewValidationError("invalid comment ID")
	}
	if err := comment.Validate(); err != nil {
		return nil, NewValidationError(err.Error())
	}

	existing, err := s.Db.GetCommentByID(comment.ID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if existing == nil {
		return nil, NewNotFoundError("comment not found")
	}
	if existing.UserID != comment.UserID {
		return nil, NewValidationError("not allowed to update comment")
	}

	if err := s.Db.UpdateComment(comment); err != nil {
		return nil, NewInternalError("error while updating comment", err)
	}

	return s.GetCommentByID(comment.ID)
}

func (s *CommentService) DeleteComment(commentID, userID int) error {
	if commentID <= 0 {
		return NewValidationError("invalid comment ID")
	}
	if userID <= 0 {
		return NewValidationError("invalid user ID")
	}

	comment, err := s.Db.GetCommentByID(commentID)
	if err != nil {
		return NewInternalError("error while querying DB", err)
	}
	if comment == nil {
		return NewNotFoundError("comment not found")
	}
	if comment.UserID != userID {
		return NewValidationError("not allowed to delete comment")
	}

	if err := s.Db.DeleteComment(commentID); err != nil {
		return NewInternalError("error while deleting comment", err)
	}
	return nil
}

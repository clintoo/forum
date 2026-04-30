package services

import "forum/backend/internal/models"

type ReactionStorer interface {
	VoteOnPost(userID, postID, voteType int) error
	VoteOnComment(userID, commentID, voteType int) error

	GetUserPostVote(userID, postID int) (int, error)
	GetUserCommentVote(userID, commentID int) (int, error)

	GetPostByID(postID int) (*models.Post, error)
	GetCommentByID(commentID int) (*models.Comment, error)
}

type ReactionService struct {
	Db ReactionStorer
}

func reactionTypeToVoteType(reactionType string) (int, error) {
	switch reactionType {
	case "upvote", "like", "up":
		return 1, nil
	case "downvote", "dislike", "down":
		return -1, nil
	default:
		return 0, NewValidationError("invalid reaction type")
	}
}

func (s *ReactionService) ReactToPost(userID, postID int, reactionType string) (*models.Post, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}
	if postID <= 0 {
		return nil, NewValidationError("invalid post ID")
	}
	if reactionType == "" {
		return nil, NewValidationError("reaction type cannot be empty")
	}

	voteType, err := reactionTypeToVoteType(reactionType)
	if err != nil {
		return nil, err
	}

	if err := s.Db.VoteOnPost(userID, postID, voteType); err != nil {
		return nil, NewInternalError("error while saving post reaction", err)
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

func (s *ReactionService) RemovePostReaction(userID, postID int) (*models.Post, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}
	if postID <= 0 {
		return nil, NewValidationError("invalid post ID")
	}

	if err := s.Db.VoteOnPost(userID, postID, 0); err != nil {
		return nil, NewInternalError("error while removing post reaction", err)
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

func (s *ReactionService) ReactToComment(userID, postID, commentID int, reactionType string) (*models.Comment, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}
	if postID <= 0 {
		return nil, NewValidationError("invalid post ID")
	}
	if commentID <= 0 {
		return nil, NewValidationError("invalid comment ID")
	}
	if reactionType == "" {
		return nil, NewValidationError("reaction type cannot be empty")
	}

	voteType, err := reactionTypeToVoteType(reactionType)
	if err != nil {
		return nil, err
	}

	if err := s.Db.VoteOnComment(userID, commentID, voteType); err != nil {
		return nil, NewInternalError("error while saving comment reaction", err)
	}

	comment, err := s.Db.GetCommentByID(commentID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if comment == nil {
		return nil, NewNotFoundError("comment not found")
	}
	return comment, nil
}

func (s *ReactionService) RemoveCommentReaction(userID, commentID int) (*models.Comment, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}
	if commentID <= 0 {
		return nil, NewValidationError("invalid comment ID")
	}

	if err := s.Db.VoteOnComment(userID, commentID, 0); err != nil {
		return nil, NewInternalError("error while removing comment reaction", err)
	}

	comment, err := s.Db.GetCommentByID(commentID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if comment == nil {
		return nil, NewNotFoundError("comment not found")
	}
	return comment, nil
}

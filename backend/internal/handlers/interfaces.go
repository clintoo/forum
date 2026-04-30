package handlers

import (
	"forum/backend/internal/dto"
	"forum/backend/internal/models"
)

// UserServicer is the interface that UserHandler and AuthHandler depend on.
type UserServicer interface {
	SignUp(user *dto.PrivateUser) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(userID int) (*models.User, error)
	GetMe(userID int) (*models.User, error)
	UpdateMe(userID int, username, email, avatarPath, password *string) (*models.User, error)
	DeleteMe(userID int) error
	GetMePosts(userID int) ([]models.Post, error)
	GetMeComments(userID int) ([]models.Comment, error)
	GetMeLikes(userID int) ([]models.Post, error)
}

// SessionServicer is the interface that AuthHandler depends on.
type SessionServicer interface {
	CreateUserSession(userID int) (*models.Session, error)
	ValidateSession(sessionID string) (int, bool)
	Logout(sessionID string) error
}

// PostServicer is the interface that ContentHandler and CategoryHandler depend on.
type PostServicer interface {
	CreatePost(post *models.Post) (*models.Post, error)
	GetAllPosts() ([]models.Post, error)
	GetPostByID(postID int) (*models.Post, error)
	UpdatePost(post *models.Post) (*models.Post, error)
	DeletePost(postID, userID int) error
	GetAllCategories() ([]models.Category, error)
	GetPostsByCategory(categoryID int) ([]models.Post, error)
}

// CommentServicer is the interface that ContentHandler depends on.
type CommentServicer interface {
	CreateComment(comment *models.Comment) (*models.Comment, error)
	GetCommentByID(commentID int) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]models.Comment, error)
	UpdateComment(comment *models.Comment) (*models.Comment, error)
	DeleteComment(commentID, userID int) error
}

// ReactionServicer is the interface that ContentHandler depends on.
type ReactionServicer interface {
	ReactToPost(userID, postID int, reactionType string) (*models.Post, error)
	RemovePostReaction(userID, postID int) (*models.Post, error)
	ReactToComment(userID, postID, commentID int, reactionType string) (*models.Comment, error)
	RemoveCommentReaction(userID, commentID int) (*models.Comment, error)
}

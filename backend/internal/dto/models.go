package dto

import (
	"fmt"
	"time"
)

// DTO is a marker interface that restricts which types may be used as response
// payloads. Every concrete DTO struct must implement IsDTO to satisfy this
// interface.
type DTO interface {
	IsDTO()
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Category) IsDTO() {}

type Comment struct {
	ID        int        `json:"id"`
	PostID    int        `json:"post_id"`
	Author    PublicUser `json:"author"`
	Content   string     `json:"content"`
	Likes     int        `json:"likes"`
	Dislikes  int        `json:"dislikes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Comment) IsDTO() {}

type Post struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	ImagePath    string     `json:"image_path,omitempty"`
	Author       PublicUser `json:"author"`
	Categories   []Category `json:"categories"`
	Likes        int        `json:"likes"`
	Dislikes     int        `json:"dislikes"`
	CommentCount int        `json:"comment_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Post) IsDTO() {}

type PostReaction struct {
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	LikeType  string    `json:"like_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (PostReaction) IsDTO() {}

type CommentReaction struct {
	UserID    int       `json:"user_id"`
	CommentID int       `json:"comment_id"`
	LikeType  string    `json:"like_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (CommentReaction) IsDTO() {}

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (Session) IsDTO() {}

type PublicUser struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	AvatarPath string    `json:"avatar_path"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PublicUser) IsDTO() {}

type PrivateUser struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	AvatarPath string    `json:"avatar_path"`
	CreatedAt  time.Time `json:"created_at"`
}

func (u PrivateUser) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	return nil
}

func (PrivateUser) IsDTO() {}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (Response) IsDTO() {}

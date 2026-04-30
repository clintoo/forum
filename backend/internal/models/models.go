package models

import (
	"fmt"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"-"` // Hide password in JSON responses
	AvatarPath string    `json:"avatar_path"`
	CreatedAt  time.Time `json:"created_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Post struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`

	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ImagePath   string   `json:"image_path,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	CategoryIDs []int    `json:"category_ids,omitempty"` // used for create/update input; not persisted directly

	Likes        int `json:"likes"`
	Dislikes     int `json:"dislikes"`
	CommentCount int `json:"comment_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) Validate() error {
	if p.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if p.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if p.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) Validate() error {
	if c.PostID <= 0 {
		return fmt.Errorf("invalid post ID")
	}
	if c.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if c.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PostStats struct {
	Likes        int
	Dislikes     int
	CommentCount int
}

type CommentStats struct {
	Likes    int
	Dislikes int
}

type PostReaction struct {
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id,omitempty"`
	LikeType  int       `json:"like_type"` // 1 = like, -1 = dislike
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentReaction struct {
	UserID    int       `json:"user_id"`
	CommentID int       `json:"comment_id,omitempty"`
	LikeType  int       `json:"like_type"` // 1 = like, -1 = dislike
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

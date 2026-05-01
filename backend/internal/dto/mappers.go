package dto

import "forum/backend/internal/models"

// ToPublicUser converts a domain User to its public transport representation.
func ToPublicUser(u *models.User) PublicUser {
	return PublicUser{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		AvatarPath: u.AvatarPath,
		CreatedAt:  u.CreatedAt,
	}
}

// ToPrivateUser converts a domain User to its private (authenticated) transport representation.
func ToPrivateUser(u *models.User) PrivateUser {
	return PrivateUser{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		AvatarPath: u.AvatarPath,
		CreatedAt:  u.CreatedAt,
	}
}

// ToCategory converts a domain Category to its transport representation.
func ToCategory(c models.Category) Category {
	return Category{ID: c.ID, Name: c.Name}
}

// ToCategories converts a slice of domain Categories to transport representations.
func ToCategories(cats []models.Category) []Category {
	result := make([]Category, len(cats))
	for i, c := range cats {
		result[i] = ToCategory(c)
	}
	return result
}

// ToPostWithAuthor converts a domain Post to its transport representation with author details.
func ToPostWithAuthor(p *models.Post, author *models.User) Post {
	cats := make([]Category, len(p.Categories))
	for i, name := range p.Categories {
		cats[i] = Category{Name: name}
	}

	authorDTO := PublicUser{ID: p.UserID}
	if author != nil {
		authorDTO = ToPublicUser(author)
	}

	return Post{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		ImagePath:    p.ImagePath,
		Author:       authorDTO,
		Categories:   cats,
		Likes:        p.Likes,
		Dislikes:     p.Dislikes,
		CommentCount: p.CommentCount,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// ToPost converts a domain Post to its transport representation.
// Categories are mapped by name only since domain queries return names, not IDs.
func ToPost(p *models.Post) Post {
	cats := make([]Category, len(p.Categories))

	for i, name := range p.Categories {
		cats[i] = Category{Name: name}
	}

	return Post{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		ImagePath:    p.ImagePath,
		Author:       PublicUser{ID: p.UserID},
		Categories:   cats,
		Likes:        p.Likes,
		Dislikes:     p.Dislikes,
		CommentCount: p.CommentCount,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// ToPosts converts a slice of domain Posts to transport representations.
func ToPosts(posts []models.Post) []Post {
	result := make([]Post, len(posts))
	for i := range posts {
		result[i] = ToPost(&posts[i])
	}
	return result
}

// ToComment converts a domain Comment to its transport representation.
func ToComment(c *models.Comment) Comment {
	return Comment{
		ID:        c.ID,
		PostID:    c.PostID,
		Author:    PublicUser{ID: c.UserID},
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToComments converts a slice of domain Comments to transport representations.
func ToComments(comments []models.Comment) []Comment {
	result := make([]Comment, len(comments))
	for i := range comments {
		result[i] = ToComment(&comments[i])
	}

	return result
}

// ToCommentWithAuthor converts a domain Comment to its transport representation with author details.
func ToCommentWithAuthor(c *models.Comment, author *models.User) Comment {
	authorDTO := PublicUser{ID: c.UserID}
	if author != nil {
		authorDTO = ToPublicUser(author)
	}

	return Comment{
		ID:        c.ID,
		PostID:    c.PostID,
		Author:    authorDTO,
		Content:   c.Content,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToPublicUsers converts a slice of domain Users to transport representations.
func ToPublicUsers(users []models.User) []PublicUser {
	result := make([]PublicUser, len(users))
	for i := range users {
		result[i] = ToPublicUser(&users[i])
	}
	return result
}

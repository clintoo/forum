package services

import (
	"forum/backend/internal/dto"
	"forum/backend/internal/models"
	"forum/backend/internal/utils"
)

type UserStorer interface {
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(userID int, username, email, avatarPath, passwordHash *string) error
	DeleteUser(userID int) error

	GetAllUsers() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserComments(userID int) ([]models.Comment, error)
	GetUserLikes(userID int) ([]models.Post, error)
	GetPostsByUserID(userID int) ([]models.Post, error)
}

type UserService struct {
	Db UserStorer
}

func (s *UserService) SignUp(user *dto.PrivateUser) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, NewValidationError(err.Error())
	}

	// - Check if email/username already exists
	userInDB, err := s.Db.GetUserByEmail(user.Email)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if userInDB != nil {
		return nil, NewConflictError("email already exists")
	}

	userInDB, err = s.Db.GetUserByUsername(user.Username)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if userInDB != nil {
		return nil, NewConflictError("username already exists")
	}

	// - Hash password
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, NewInternalError("error while hashing password", err)
	}

	// - Create user in database
	created, err := s.Db.CreateUser(&models.User{
		Email:      user.Email,
		Username:   user.Username,
		AvatarPath: user.AvatarPath,
		Password:   hashed,
	})
	if err != nil {
		return nil, NewInternalError("error while creating user", err)
	}

	return created, nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	// Validate required fields
	if email == "" {
		return nil, NewValidationError("email cannot be empty")
	}
	if password == "" {
		return nil, NewValidationError("password cannot be empty")
	}

	// Get user from database by email
	user, err := s.Db.GetUserByEmail(email)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}

	// Check if user exists
	if user == nil {
		return nil, NewValidationError("user not found")
	}

	// Compare password with hash using the checkPasswordHash function from utils
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, NewValidationError("invalid password")
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.Db.GetAllUsers()
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return users, nil
}

func (s *UserService) GetUserByID(userID int) (*models.User, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}

	user, err := s.Db.GetUserByID(userID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	if user == nil {
		return nil, NewNotFoundError("user not found")
	}
	return user, nil
}

func (s *UserService) GetMe(userID int) (*models.User, error) {
	return s.GetUserByID(userID)
}

func (s *UserService) UpdateMe(userID int, username, email, avatarPath, password *string) (*models.User, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}

	if email != nil {
		if *email == "" {
			return nil, NewValidationError("email cannot be empty")
		}
		existing, err := s.Db.GetUserByEmail(*email)
		if err != nil {
			return nil, NewInternalError("error while querying DB", err)
		}
		if existing != nil && existing.ID != userID {
			return nil, NewConflictError("email already exists")
		}
	}

	if username != nil {
		if *username == "" {
			return nil, NewValidationError("username cannot be empty")
		}
		existing, err := s.Db.GetUserByUsername(*username)
		if err != nil {
			return nil, NewInternalError("error while querying DB", err)
		}
		if existing != nil && existing.ID != userID {
			return nil, NewConflictError("username already exists")
		}
	}

	var passwordHash *string
	if password != nil {
		if *password == "" {
			return nil, NewValidationError("password cannot be empty")
		}
		hashed, err := utils.HashPassword(*password)
		if err != nil {
			return nil, NewInternalError("error while hashing password", err)
		}
		passwordHash = &hashed
	}

	if err := s.Db.UpdateUser(userID, username, email, avatarPath, passwordHash); err != nil {
		return nil, NewInternalError("error while updating user", err)
	}

	return s.GetUserByID(userID)
}

func (s *UserService) DeleteMe(userID int) error {
	if userID <= 0 {
		return NewValidationError("invalid user ID")
	}

	if err := s.Db.DeleteUser(userID); err != nil {
		return NewInternalError("error while deleting user", err)
	}
	return nil
}

func (s *UserService) GetMePosts(userID int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}

	posts, err := s.Db.GetPostsByUserID(userID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return posts, nil
}

func (s *UserService) GetMeComments(userID int) ([]models.Comment, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}

	comments, err := s.Db.GetUserComments(userID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return comments, nil
}

func (s *UserService) GetMeLikes(userID int) ([]models.Post, error) {
	if userID <= 0 {
		return nil, NewValidationError("invalid user ID")
	}

	likes, err := s.Db.GetUserLikes(userID)
	if err != nil {
		return nil, NewInternalError("error while querying DB", err)
	}
	return likes, nil
}

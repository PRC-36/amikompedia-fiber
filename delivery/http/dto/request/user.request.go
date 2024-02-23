package request

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
)

// UserRequest to create
type UserRequest struct {
	Email    string `json:"email"`
	Nim      string `json:"nim"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserResetPasswordRequest struct {
	Password        string `json:"password" validate:"required,min=8,containsany,containsuppercase,containslowercase,containsnumeric"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserUpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,containsany,containsuppercase,containslowercase,containsnumeric"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UserFollowRequest struct {
	FollowID string `json:"follow_id" validate:"required"`
}

func (u UserRequest) ToEntity() *entity.User {
	return &entity.User{
		Email:    u.Email,
		Nim:      u.Nim,
		Name:     u.Name,
		Username: u.Username,
		Password: u.Password,
		Bio:      u.Bio,
	}
}

// UserUpdateRequest to update
type UserUpdateRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio"`
}

func (u UserUpdateRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: u.Username,
		Name:     u.Name,
		Bio:      u.Bio,
	}
}

func (u UserUpdatePasswordRequest) ToUserEntity(userId sql.NullString) *entity.User {
	return &entity.User{
		ID:       userId,
		Password: u.NewPassword,
	}
}

func (u UserFollowRequest) ToUserFollowEntity(userId string) *entity.UserFollow {
	return &entity.UserFollow{
		FollowerID:  userId,
		FollowingID: u.FollowID,
	}
}

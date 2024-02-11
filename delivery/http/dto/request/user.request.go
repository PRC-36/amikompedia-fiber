package request

import (
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

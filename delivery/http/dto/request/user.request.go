package request

import (
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"github.com/google/uuid"
)

type UserRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,number"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (u UserRegisterRequest) ToEntity() *entity.User {
	return &entity.User{
		ID:        uuid.NewString(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Password:  u.Password,
	}
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u UserLoginRequest) ToEntity() *entity.User {
	return &entity.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

package request

import "github.com/PRC-36/amikompedia-fiber/domain/entity"

type UserRegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Nim             string `json:"nim" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,containsany,containsuppercase,containslowercase,containsnumeric"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (r UserRegisterRequest) ToEntity() *entity.UserRegister {
	return &entity.UserRegister{
		Name:     r.Name,
		Email:    r.Email,
		Nim:      r.Nim,
		Password: r.Password,
	}
}

package request

import "github.com/PRC-36/amikompedia-fiber/domain/entity"

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

func (r LoginRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: r.UsernameOrEmail,
		Email:    r.UsernameOrEmail,
		Password: r.Password,
	}
}

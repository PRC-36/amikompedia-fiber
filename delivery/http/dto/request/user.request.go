package request

import (
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
)

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

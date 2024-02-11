package entity

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"time"
)

type UserRegister struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement;not null"`
	Name            string    `gorm:"column:name"`
	Email           string    `gorm:"column:email"`
	Nim             string    `gorm:"column:nim"`
	Password        string    `gorm:"column:password"`
	IsVerified      bool      `gorm:"column:is_verified"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (e *UserRegister) TableName() string {
	return "user_register"
}

func (e *UserRegister) ToUserRegisterResponse(refCode string) *response.RegisterResponseWithRefCode {
	return &response.RegisterResponseWithRefCode{
		RefCode: refCode,
		RegisterUser: response.RegisterResponse{
			ID:         e.ID,
			Name:       e.Name,
			Email:      e.Email,
			Nim:        e.Nim,
			IsVerified: e.IsVerified,
			CreatedAt:  e.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}
}

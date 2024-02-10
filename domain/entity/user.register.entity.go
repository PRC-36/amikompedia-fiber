package entity

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"time"
)

type UserRegister struct {
	ID              int       `gorm:"primaryKey;autoIncrement;not null"`
	Name            string    `gorm:"type:varchar(255);not null"`
	Email           string    `gorm:"type:varchar(255);not null"`
	Nim             string    `gorm:"type:varchar(10);not null"`
	Password        string    `gorm:"type:varchar(255);not null"`
	IsVerified      bool      `gorm:"type:boolean;not null;default:false"`
	EmailVerifiedAt time.Time `gorm:"type:timestamp;null"`
	CreatedAt       time.Time `gorm:"type:timestamp;not null"`
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

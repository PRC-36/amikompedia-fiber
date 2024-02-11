package entity

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"time"
)

type User struct {
	ID        sql.NullString `gorm:"column:uuid;primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email     string         `gorm:"column:email"`
	Nim       string         `gorm:"column:nim"`
	Name      string         `gorm:"column:name"`
	Username  string         `gorm:"column:username"`
	Bio       string         `gorm:"column:bio"`
	Password  string         `gorm:"column:password"`
	Images    []Image        `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

func (e *User) TableName() string {
	return "users"
}

func (u *User) ToUserResponse() *response.UserResponse {

	return &response.UserResponse{
		ID:        u.ID.String,
		Email:     u.Email,
		Nim:       u.Nim,
		Name:      u.Name,
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     ToImageResponses(u.Images),
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (u *User) ToUserResponseWithToken(sessionResp *response.SessionsResponse) *response.LoginResponse {
	return &response.LoginResponse{
		Token: sessionResp,
		User:  u.ToUserResponse(),
	}
}

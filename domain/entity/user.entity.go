package entity

import "github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"

type User struct {
	ID        string `gorm:"column:UserID;primaryKey"`
	FirstName string `gorm:"column:FirstName"`
	LastName  string `gorm:"column:LastName"`
	Email     string `gorm:"column:Email"`
	Phone     string `gorm:"column:PhoneNumber"`
	Password  string `gorm:"column:Password"`
}

func (u *User) TableName() string {
	return "Users"
}

func (u *User) ToUserResponse() *response.UserResponse {
	return &response.UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
	}
}

func (u *User) ToUserLoginResponse(token string) *response.UserLoginResponse {
	return &response.UserLoginResponse{
		Token: token,
		User:  *u.ToUserResponse(),
	}
}

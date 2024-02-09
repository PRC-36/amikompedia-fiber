package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	Create(DB *gorm.DB, value *entity.User) error
	CheckEmailIsExists(DB *gorm.DB, email string) (bool, *entity.User, error)
}

type userRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (u *userRepositoryImpl) Create(DB *gorm.DB, value *entity.User) error {

	result := u.DB.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when creating user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) CheckEmailIsExists(DB *gorm.DB, email string) (bool, *entity.User, error) {

	var user entity.User

	result := u.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when checking email : %v", result.Error))
		return false, &user, result.Error
	}

	return true, &user, nil
}

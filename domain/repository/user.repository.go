package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	UserCreate(tx *gorm.DB, value *entity.User) error
	FindByUsernameOrEmail(tx *gorm.DB, value *entity.User) error
}

type userRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (u *userRepositoryImpl) UserCreate(tx *gorm.DB, value *entity.User) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) FindByUsernameOrEmail(tx *gorm.DB, value *entity.User) error {
	result := tx.Where("username = ? OR email = ?", value.Username, value.Email).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find username or email : %v", result.Error))
		return result.Error
	}

	return nil
}

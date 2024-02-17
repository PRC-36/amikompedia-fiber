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
	FindByUserUUID(tx *gorm.DB, value *entity.User) error
	UpdateUser(tx *gorm.DB, value *entity.User) error
	UpdatePassword(tx *gorm.DB, value *entity.User) error
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

func (u *userRepositoryImpl) FindByUserUUID(tx *gorm.DB, value *entity.User) error {
	result := tx.Preload("Images").Where("uuid = ?", value.ID).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user by user uuid : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) UpdateUser(tx *gorm.DB, value *entity.User) error {
	result := tx.Model(value).Updates(value).Preload("Images").First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when update user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) UpdatePassword(tx *gorm.DB, value *entity.User) error {
	result := tx.Model(value).Where("uuid = ?", value.ID).Update("password", value.Password)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when set new password : %v", result.Error))
		panic(result.Error)
	}

	return nil
}

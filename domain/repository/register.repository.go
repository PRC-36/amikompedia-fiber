package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type RegisterRepository interface {
	RegisterCreate(tx *gorm.DB, value *entity.UserRegister) error
	RegisterCheckEmailIsVerified(tx *gorm.DB, email string) (bool, *entity.UserRegister, error)
	FindById(tx *gorm.DB, id int) (*entity.UserRegister, error)
	Update(tx *gorm.DB, value *entity.UserRegister) *entity.UserRegister
}

type registerRepositoryImpl struct {
}

func NewRegisterRepository() RegisterRepository {
	return &registerRepositoryImpl{}
}

func (r *registerRepositoryImpl) RegisterCreate(tx *gorm.DB, value *entity.UserRegister) error {

	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when creating user register : %v", result.Error))
		return result.Error
	}

	return nil
}

func (r *registerRepositoryImpl) RegisterCheckEmailIsVerified(tx *gorm.DB, email string) (bool, *entity.UserRegister, error) {
	var user entity.UserRegister

	result := tx.Where("email = ?", email).First(&user)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when checking email : %v", result.Error))
		return false, nil, result.Error
	}

	return user.IsVerified, &user, nil
}

func (r *registerRepositoryImpl) FindById(tx *gorm.DB, id int) (*entity.UserRegister, error) {
	var user entity.UserRegister

	result := tx.Where("id = ?", id).First(&user)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user by id : %v", result.Error))
		return nil, result.Error
	}

	return &user, nil
}

func (r *registerRepositoryImpl) Update(tx *gorm.DB, value *entity.UserRegister) *entity.UserRegister {
	result := tx.Model(&value).Where("id = ?", value.ID).Updates(entity.UserRegister{IsVerified: value.IsVerified, EmailVerifiedAt: value.EmailVerifiedAt})

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when updating user register : %v", result.Error))
		return nil
	}

	return value
}

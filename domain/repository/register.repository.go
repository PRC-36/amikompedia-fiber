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
	var user *entity.UserRegister

	result := tx.Where("email = ?", email).First(user)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when checking email : %v", result.Error))
		return false, user, result.Error
	}

	return user.IsVerified, user, nil
}

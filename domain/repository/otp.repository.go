package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type OtpRepository interface {
	OtpCreate(tx *gorm.DB, value *entity.Otp) error
}

type otpRepositoryImpl struct {
}

func NewOtpRepository() OtpRepository {
	return &otpRepositoryImpl{}
}

func (o *otpRepositoryImpl) OtpCreate(tx *gorm.DB, value *entity.Otp) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when creating otp : %v", result.Error))
		return result.Error
	}

	return nil
}

package repository

import (
	"errors"
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type OtpRepository interface {
	OtpCreate(tx *gorm.DB, value *entity.Otp) error
	OtpUpdate(tx *gorm.DB, value *entity.Otp) error
	FindByRefCode(tx *gorm.DB, refCode string) (*entity.Otp, error)
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

func (o *otpRepositoryImpl) OtpUpdate(tx *gorm.DB, value *entity.Otp) error {
	result := tx.Model(&value).Where("ref_code = ?", value.RefCode).Updates(entity.Otp{OtpValue: value.OtpValue, ExpiredAt: value.ExpiredAt})
	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when updating otp : %v", result.Error))
		return result.Error
	}

	return nil
}

func (o *otpRepositoryImpl) FindByRefCode(tx *gorm.DB, refCode string) (*entity.Otp, error) {
	var otp entity.Otp
	result := tx.Preload("UserRegister").Preload("User").Where("ref_code = ?", refCode).First(&otp)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refferal code not found")
		}
		return nil, result.Error
	}
	return &otp, nil
}

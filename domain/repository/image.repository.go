package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type ImageRepository interface {
	ImageSave(tx *gorm.DB, value *entity.Image) error
	ImageFindByUserID(tx *gorm.DB, value *entity.Image) error
}

type imageRepositoryImpl struct {
}

func NewImageRepository() ImageRepository {
	return &imageRepositoryImpl{}
}

func (i *imageRepositoryImpl) ImageSave(tx *gorm.DB, value *entity.Image) error {
	result := tx.Save(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when saving image : %v", result.Error))
		return result.Error
	}

	return nil
}

func (i *imageRepositoryImpl) ImageFindByUserID(tx *gorm.DB, value *entity.Image) error {
	result := tx.Where("user_uuid = ?", value.UserID).Where("image_type = ?", value.ImageType).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find image by user id : %v", result.Error))
		return result.Error
	}

	return nil
}

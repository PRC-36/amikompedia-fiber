package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type UserFollowRepository interface {
	UserFollowCreate(tx *gorm.DB, value *entity.UserFollow) error
	UserFollowDelete(tx *gorm.DB, value *entity.UserFollow) error
	FindByFollowID(tx *gorm.DB, value *entity.UserFollow) (*entity.UserFollow, error)
}

type userFollowRepositoryImpl struct {
}

func NewUserFollowRepository() UserFollowRepository {
	return &userFollowRepositoryImpl{}
}

func (u *userFollowRepositoryImpl) UserFollowCreate(tx *gorm.DB, value *entity.UserFollow) error {

	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create follow : %v", result.Error))
		return result.Error
	}

	return nil

}

func (u *userFollowRepositoryImpl) UserFollowDelete(tx *gorm.DB, value *entity.UserFollow) error {

	result := tx.Where("follower_id = ? AND following_id = ?", value.FollowerID, value.FollowingID).Delete(&entity.UserFollow{})

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when delete follow : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userFollowRepositoryImpl) FindByFollowID(tx *gorm.DB, value *entity.UserFollow) (*entity.UserFollow, error) {

	result := tx.Where("follower_id = ? AND following_id = ?", value.FollowerID, value.FollowingID).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find follow by id : %v", result.Error))
		return nil, result.Error
	}

	return value, nil
}

package usecase

import (
	"context"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type UserUsecase interface {
	CreateNewUser(ctx context.Context, requestData *request.UserRequest) (*response.UserResponse, error)
}

type userUsecaseImpl struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	UserRepository repository.UserRepository
}

func NewUserUsecase(db *gorm.DB, validate *validator.Validate, userRepository repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{DB: db, Validate: validate, UserRepository: userRepository}
}

func (u *userUsecaseImpl) CreateNewUser(ctx context.Context, requestData *request.UserRequest) (*response.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	hashedPassword, _ := util.HashPassword(requestData.Password)
	requestData.Password = hashedPassword

	requestUserEntity := requestData.ToEntity()
	err = u.UserRepository.UserCreate(tx, requestUserEntity)

	if err != nil {
		log.Printf("Failed create user : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return requestUserEntity.ToUserResponse(), nil
}

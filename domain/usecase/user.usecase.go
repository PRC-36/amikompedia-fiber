package usecase

import (
	"context"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"log"
)

type UserUsecase interface {
	Register(ctx context.Context, request *request.UserRegisterRequest) (*response.UserResponse, error)
	Login(ctx context.Context, request *request.UserLoginRequest) (*response.UserLoginResponse, error)
}

type userUsecaseImpl struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	UserRepository repository.UserRepository
	TokenMaker     token.Maker
}

func NewUserUsecase(db *gorm.DB, validate *validator.Validate, userRepository repository.UserRepository, tokenMaker token.Maker) UserUsecase {
	return &userUsecaseImpl{DB: db, Validate: validate, UserRepository: userRepository, TokenMaker: tokenMaker}
}

func (u *userUsecaseImpl) Register(ctx context.Context, request *request.UserRegisterRequest) (*response.UserResponse, error) {

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(request)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	isExists, _, _ := u.UserRepository.CheckEmailIsExists(tx, request.Email)
	if isExists {
		log.Printf("Email already exists")
		return nil, util.EmailAlreadyExist
	}

	hashedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		log.Printf("Failed hash password : %+v", err)
		return nil, err
	}

	request.Password = hashedPassword
	requestEntity := request.ToEntity()

	err = u.UserRepository.Create(tx, requestEntity)

	if err != nil {
		log.Printf("Failed create user : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return requestEntity.ToUserResponse(), nil
}

func (u *userUsecaseImpl) Login(ctx context.Context, request *request.UserLoginRequest) (*response.UserLoginResponse, error) {

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(request)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	isExists, user, _ := u.UserRepository.CheckEmailIsExists(tx, request.Email)
	if !isExists {
		log.Printf("Email not found")
		return nil, util.EmailNotFound
	}

	if !util.CheckPassword(request.Password, user.Password) {
		log.Printf("Invalid password")
		return nil, util.InvalidPassword
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed commit transaction : %+v", err)
		return nil, err
	}

	token, err := u.TokenMaker.CreateToken(token.UserPayload{UUID: user.ID, Username: user.Email})

	if err != nil {
		log.Printf("Failed create token : %+v", err)
		return nil, err
	}

	return user.ToUserLoginResponse(token), nil

}

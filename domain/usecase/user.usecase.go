package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/aws"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
)

type UserUsecase interface {
	CreateNewUser(ctx context.Context, requestData *request.UserRequest) (*response.UserResponse, error)
	ProfileUser(ctx context.Context, userID string) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, userID string, requestData *request.UserUpdateRequest, imgAvtr, imgHeader *multipart.FileHeader) (*response.UserResponse, error)
}

type userUsecaseImpl struct {
	DB              *gorm.DB
	Validate        *validator.Validate
	AwsS3           aws.AwsS3Action
	UserRepository  repository.UserRepository
	ImageRepository repository.ImageRepository
}

func NewUserUsecase(db *gorm.DB, validate *validator.Validate, awsS3 aws.AwsS3Action,
	userRepository repository.UserRepository, imageRepository repository.ImageRepository) UserUsecase {
	return &userUsecaseImpl{DB: db, Validate: validate,
		AwsS3: awsS3, UserRepository: userRepository, ImageRepository: imageRepository}
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

func (u *userUsecaseImpl) ProfileUser(ctx context.Context, userID string) (*response.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Var(userID, "required")
	if err != nil {
		log.Printf("Invalid username : %+v", err)
		return nil, err
	}

	userEntity := entity.User{ID: sql.NullString{Valid: true, String: userID}}

	err = u.UserRepository.FindByUserUUID(tx, &userEntity)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Username or email not found : %+v", err)
			return nil, util.UsernameOrEmailNotFound
		}
		log.Printf("Failed find username or email : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return userEntity.ToUserResponse(), nil
}

func (u *userUsecaseImpl) UpdateUser(ctx context.Context, userID string, requestData *request.UserUpdateRequest, imgAvtr, imgHeader *multipart.FileHeader) (*response.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	err = u.UserRepository.FindByUserUUID(tx, &entity.User{ID: sql.NullString{Valid: true, String: userID}})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found : %+v", err)
			return nil, util.UserNotFound
		}
		log.Printf("Failed find user : %+v", err)
		return nil, err
	}

	if imgAvtr != nil {

		img1, err := u.AwsS3.UploadFile(imgAvtr, aws.ImgAvatar)
		if err != nil {
			log.Printf("Failed upload avatar image : %+v", err)
			return nil, err
		}

		imgAvtrEntity := entity.Image{UserID: userID, ImageType: img1.ImageType}
		err = u.ImageRepository.ImageFindByUserID(tx, &imgAvtrEntity)
		if err != nil {
			log.Printf("Failed find image avatar : %+v", err)
		}

		imgAvtrEntity.ImageUrl = img1.ImageUrl
		imgAvtrEntity.FilePath = img1.FilePath

		err = u.ImageRepository.ImageSave(tx, &imgAvtrEntity)
		if err != nil {
			log.Printf("Failed upload save image avatar : %+v", err)
			return nil, err
		}
	}

	if imgHeader != nil {

		img2, err := u.AwsS3.UploadFile(imgHeader, aws.ImgHeader)
		if err != nil {
			log.Printf("Failed upload header image : %+v", err)
			return nil, err
		}

		imgHeaderEntity := entity.Image{UserID: userID, ImageType: img2.ImageType}
		err = u.ImageRepository.ImageFindByUserID(tx, &imgHeaderEntity)
		if err != nil {
			log.Printf("Failed find image header : %+v", err)
		}

		imgHeaderEntity.ImageUrl = img2.ImageUrl
		imgHeaderEntity.FilePath = img2.FilePath
		err = u.ImageRepository.ImageSave(tx, &imgHeaderEntity)
		if err != nil {
			log.Printf("Failed upload save image header : %+v", err)
			return nil, err
		}
	}

	userEntity := requestData.ToEntity()
	userEntity.ID = sql.NullString{Valid: true, String: userID}
	err = u.UserRepository.UpdateUser(tx, userEntity)

	if err != nil {
		log.Printf("Failed update user : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return userEntity.ToUserResponse(), nil
}

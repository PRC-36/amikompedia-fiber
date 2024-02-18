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
	"github.com/PRC-36/amikompedia-fiber/shared/mail"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"strconv"
	"time"
)

type UserUsecase interface {
	CreateNewUser(ctx context.Context, requestData *request.UserRequest) (*response.UserResponse, error)
	ProfileUser(ctx context.Context, userID string) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, userID string, requestData *request.UserUpdateRequest, imgAvtr, imgHeader *multipart.FileHeader) (*response.UserResponse, error)
	ForgotPassword(ctx context.Context, requestData *request.UserForgotPasswordRequest) (*response.OtpResponse, error)
	ResetPassword(ctx context.Context, requestData *request.UserResetPasswordRequest, secretToken string) error
	UpdatePassword(ctx context.Context, userID string, requestData *request.UserUpdatePasswordRequest) error
}

type userUsecaseImpl struct {
	DB              *gorm.DB
	Validate        *validator.Validate
	AwsS3           aws.AwsS3Action
	EmailSender     mail.EmailSender
	TokenMaker      token.Maker
	ViperConfig     util.Config
	UserRepository  repository.UserRepository
	ImageRepository repository.ImageRepository
	OtpRepository   repository.OtpRepository
}

func NewUserUsecase(DB *gorm.DB, validate *validator.Validate, awsS3 aws.AwsS3Action, emailSender mail.EmailSender, tokenMaker token.Maker, viperConfig util.Config, userRepository repository.UserRepository, imageRepository repository.ImageRepository, otpRepository repository.OtpRepository) UserUsecase {
	return &userUsecaseImpl{DB: DB, Validate: validate, AwsS3: awsS3, EmailSender: emailSender, TokenMaker: tokenMaker, ViperConfig: viperConfig, UserRepository: userRepository, ImageRepository: imageRepository, OtpRepository: otpRepository}
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

func (u *userUsecaseImpl) ForgotPassword(ctx context.Context, requestData *request.UserForgotPasswordRequest) (*response.OtpResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(requestData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Invalid email : %+v", err)
			return nil, util.EmailNotFound
		}
		log.Printf("Invalid email : %+v", err)
		return nil, err
	}

	userEntity := entity.User{Email: requestData.Email}

	err = u.UserRepository.FindByUsernameOrEmail(tx, &userEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Email not found : %+v", err)
			return nil, util.EmailNotFound
		}
		log.Printf("Failed find email : %+v", err)
		return nil, err
	}

	createRequestOtp := request.OtpCreateRequest{
		OtpValue:  strconv.FormatInt(util.RandomInt(100000, 999999), 10),
		RefCode:   util.RandomCombineIntAndString(),
		ExpiredAt: time.Now().Add(time.Minute * 1),
		UserID:    userEntity.ID,
	}

	createRequestEntity := createRequestOtp.ToEntity()

	err = u.OtpRepository.OtpCreate(tx, createRequestEntity)
	if err != nil {
		log.Printf("Failed create otp  : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	go func() {
		subject, content, toEmail := mail.GetSenderParamEmailRegist(requestData.Email, createRequestEntity.OtpValue)
		err := u.EmailSender.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
		if err != nil {
			log.Printf("Failed send email : %+v", err)
		}
	}()

	return createRequestEntity.ToOtpResponse(), nil
}

func (u *userUsecaseImpl) ResetPassword(ctx context.Context, requestData *request.UserResetPasswordRequest, secretToken string) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	log.Printf(secretToken)
	payload, err := u.TokenMaker.VerifyToken(secretToken, u.ViperConfig.TokenResetPasswordKey)
	if err != nil {

		log.Printf("Failed verify token : %+v", err)
		return err
	}

	err = u.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return err
	}

	userEntity := &entity.User{
		ID: sql.NullString{Valid: true, String: payload.UserID},
	}

	err = u.UserRepository.FindByUserUUID(tx, userEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found : %+v", err)
			return util.UserNotFound
		}
		log.Printf("Failed find user : %+v", err)
		return err
	}

	password, err := util.HashPassword(requestData.Password)
	if err != nil {
		log.Printf("Failed hash password : %+v", err)
		return err
	}

	userEntity.Password = password

	err = u.UserRepository.UpdatePassword(tx, userEntity)
	if err != nil {
		log.Printf("Failed update password : %+v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) UpdatePassword(ctx context.Context, userId string, requestData *request.UserUpdatePasswordRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return err
	}

	userEntity := requestData.ToUserEntity(sql.NullString{Valid: true, String: userId})

	err = u.UserRepository.FindByUserUUID(tx, userEntity)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found : %+v", err)
			return util.UserNotFound
		}
		log.Printf("Failed find user : %+v", err)
		return err
	}

	if !util.CheckPassword(requestData.CurrentPassword, userEntity.Password) {
		return util.CurrentPasswordDoesNotMatch
	}

	password, err := util.HashPassword(requestData.NewPassword)
	if err != nil {
		log.Printf("Failed hash password : %+v", err)
		return err
	}

	userEntity.Password = password

	err = u.UserRepository.UpdatePassword(tx, userEntity)
	if err != nil {
		log.Printf("Failed update password : %+v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

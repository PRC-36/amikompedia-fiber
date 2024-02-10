package usecase

import (
	"context"
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/mail"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type LoginUsecase interface {
	Login(ctx context.Context, requestData *request.LoginRequest) (*response.LoginResponse, error)
}

type loginUsecase struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	EmailSender    mail.EmailSender
	TokenMaker     token.Maker
	UserRepository repository.UserRepository
}

func NewLoginUsecase(db *gorm.DB, validate *validator.Validate, emailSender mail.EmailSender, tokenMaker token.Maker, userRepository repository.UserRepository) LoginUsecase {
	return &loginUsecase{DB: db,
		Validate:       validate,
		EmailSender:    emailSender,
		TokenMaker:     tokenMaker,
		UserRepository: userRepository,
	}
}

func (l *loginUsecase) Login(ctx context.Context, requestData *request.LoginRequest) (*response.LoginResponse, error) {
	tx := l.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := l.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	requestUserEntity := requestData.ToEntity()

	err = l.UserRepository.FindByUsernameOrEmail(tx, requestUserEntity)

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

	isValid := util.CheckPassword(requestData.Password, requestUserEntity.Password)
	if !isValid {
		log.Printf("Password not valid")
		return nil, util.InvalidPassword
	}

	token, err := l.TokenMaker.CreateToken(token.UserPayload{
		UserID:   requestUserEntity.ID.String,
		Username: requestUserEntity.Username,
	})

	if err != nil {
		log.Printf("Failed create token : %+v", err)
		return nil, err
	}

	return requestUserEntity.ToUserResponseWithToken(token), nil
}

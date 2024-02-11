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
	Login(ctx context.Context, requestData *request.LoginRequest, userAgent, clientIP string) (*response.LoginResponse, error)
}

type loginUsecase struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	EmailSender       mail.EmailSender
	TokenMaker        token.Maker
	ViperConfig       util.Config
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
}

func NewLoginUsecase(db *gorm.DB, validate *validator.Validate, emailSender mail.EmailSender, tokenMaker token.Maker, viperConfig util.Config, userRepository repository.UserRepository, sessionRepository repository.SessionRepository) LoginUsecase {
	return &loginUsecase{DB: db,
		Validate:          validate,
		EmailSender:       emailSender,
		TokenMaker:        tokenMaker,
		ViperConfig:       viperConfig,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
	}
}

func (l *loginUsecase) Login(ctx context.Context, requestData *request.LoginRequest, userAgent, clientIP string) (*response.LoginResponse, error) {
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

	isValid := util.CheckPassword(requestData.Password, requestUserEntity.Password)
	if !isValid {
		log.Printf("Password not valid")
		return nil, util.InvalidPassword
	}

	accessToken, accessPayload, err := l.TokenMaker.CreateToken(token.UserPayload{
		UserID:   requestUserEntity.ID.String,
		Username: requestUserEntity.Username,
	}, l.ViperConfig.TokenAccessSymetricKey, l.ViperConfig.TokenAccessDuration)

	if err != nil {
		log.Printf("Failed create access token : %+v", err)
		return nil, err
	}

	refreshToken, refreshPayload, err := l.TokenMaker.CreateToken(token.UserPayload{
		UserID:   requestUserEntity.ID.String,
		Username: requestUserEntity.Username,
	}, l.ViperConfig.TokenRefreshSymetricKey, l.ViperConfig.RefreshTokenDuration)

	if err != nil {
		log.Printf("Failed create refresh token : %+v", err)
		return nil, err
	}

	sessionRequest := &request.SessionRequest{
		ID:           refreshPayload.ID.String(),
		UserID:       requestUserEntity.ID.String,
		Username:     requestUserEntity.Username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	sessionEntity := sessionRequest.ToEntity()

	err = l.SessionRepository.SessionCreate(tx, sessionEntity)

	if err != nil {
		log.Printf("Failed create session user : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	sessionResponse := &response.SessionsResponse{
		SessionsID:            sessionEntity.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt.Format("2006-01-02 15:04:05"),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt.Format("2006-01-02 15:04:05"),
	}
	return requestUserEntity.ToUserResponseWithToken(sessionResponse), nil
}

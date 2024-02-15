package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/mail"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"time"
)

type OtpUsecase interface {
	OtpCreate(ctx context.Context, requestData *request.OtpCreateRequest) (*response.OtpResponse, error)
	OtpValidation(ctx context.Context, requestData *request.OtpValidateRequest, userAgent, clientIP string) (*response.LoginResponse, error)
	ResendOtp(ctx context.Context, requestData *request.OtpSendRequest) error
}

type otpUsecaseImpl struct {
	DB                 *gorm.DB
	Validate           *validator.Validate
	EmailSender        mail.EmailSender
	TokenMaker         token.Maker
	ViperConfig        util.Config
	OtpRepository      repository.OtpRepository
	RegisterRepository repository.RegisterRepository
	UserRepository     repository.UserRepository
	SessionRepository  repository.SessionRepository
}

func NewOtpUsecase(DB *gorm.DB, validate *validator.Validate, emailSender mail.EmailSender, tokenMaker token.Maker, viperConfig util.Config, otpRepository repository.OtpRepository, registerRepository repository.RegisterRepository, userRepository repository.UserRepository, sessionRepository repository.SessionRepository) OtpUsecase {
	return &otpUsecaseImpl{DB: DB, Validate: validate, EmailSender: emailSender, TokenMaker: tokenMaker, ViperConfig: viperConfig, OtpRepository: otpRepository, RegisterRepository: registerRepository, UserRepository: userRepository, SessionRepository: sessionRepository}
}

func (o otpUsecaseImpl) OtpCreate(ctx context.Context, requestData *request.OtpCreateRequest) (*response.OtpResponse, error) {
	tx := o.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := o.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	requestOtpEntity := requestData.ToEntity()

	err = o.OtpRepository.OtpCreate(tx, requestOtpEntity)
	if err != nil {
		log.Printf("Failed create otp: %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return requestOtpEntity.ToOtpResponse(), nil
}

func (o otpUsecaseImpl) OtpValidation(ctx context.Context, requestData *request.OtpValidateRequest, userAgent, clientIP string) (*response.LoginResponse, error) {
	tx := o.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := o.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	result, err := o.OtpRepository.FindByRefCode(tx, requestData.RefCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Failed find otp by ref code : %+v", err)
			return nil, util.RefCodeNotFound
		}
		log.Printf("Failed find otp by ref code : %+v", err)
		return nil, err
	}

	if result.ExpiredAt.Before(time.Now()) {
		log.Printf("OTP Expired")
		return nil, fmt.Errorf("OTP Expired. Please requesy a new one")
	}

	if result.OtpValue != requestData.OtpValue {
		log.Printf("OTP Value not match")
		return nil, fmt.Errorf("OTP Value not match")
	}

	if result.UserRID.Valid {
		userRid, err := o.RegisterRepository.FindById(tx, int(result.UserRID.Int32))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Failed find user register by id : %+v", err)
				return nil, util.UserRegisterNotFound
			}
			log.Printf("Failed find user register by id : %+v", err)
			return nil, err
		}

		userName := util.RandomString(7)
		userEntity := userRid.ToUserEntity(userName)

		paramsUpdate := entity.UserRegister{
			IsVerified:      true,
			EmailVerifiedAt: time.Now(),
			ID:              int(result.UserRID.Int32),
		}
		o.RegisterRepository.Update(tx, &paramsUpdate)

		err = o.UserRepository.UserCreate(tx, userEntity)
		if err != nil {
			log.Printf("Failed create user : %+v", err)
			return nil, err
		}

		accessToken, accessPayload, err := o.TokenMaker.CreateToken(token.UserPayload{
			UserID:   userEntity.ID.String,
			Username: userEntity.Username,
		}, o.ViperConfig.TokenAccessSymetricKey, o.ViperConfig.TokenAccessDuration)

		if err != nil {
			log.Printf("Failed create refresh token : %+v", err)
			return nil, err
		}

		refreshToken, refreshPayload, err := o.TokenMaker.CreateToken(token.UserPayload{
			UserID:   userEntity.ID.String,
			Username: userEntity.Username,
		}, o.ViperConfig.TokenRefreshSymetricKey, o.ViperConfig.RefreshTokenDuration)

		if err != nil {
			log.Printf("Failed create refresh token : %+v", err)
			return nil, err
		}

		sessionRequest := &request.SessionRequest{
			ID:           refreshPayload.ID.String(),
			UserID:       userEntity.ID.String,
			Username:     userEntity.Username,
			RefreshToken: refreshToken,
			UserAgent:    userAgent,
			ClientIP:     clientIP,
			IsBlocked:    false,
			ExpiredAt:    refreshPayload.ExpiredAt,
		}

		sessionEntity := sessionRequest.ToEntity()

		err = o.SessionRepository.SessionCreate(tx, sessionEntity)

		if err != nil {
			log.Printf("Failed create session : %+v", err)
			return nil, err
		}

		err = tx.Commit().Error
		if err != nil {
			return nil, err
		}

		sessionResponse := &response.SessionsResponse{
			SessionsID:            sessionEntity.ID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		}

		return userEntity.ToUserResponseWithToken(sessionResponse), nil
	}
	if result.UserID.Valid {
		return nil, nil
	}

	return nil, err
}

func (o otpUsecaseImpl) ResendOtp(ctx context.Context, requestData *request.OtpSendRequest) error {
	tx := o.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newOtp := util.GenerateOtpValue()

	result, err := o.OtpRepository.FindByRefCode(tx, requestData.RefCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Failed find otp by ref code : %+v", err)
			return util.RefCodeNotFound
		}
		log.Printf("Failed find otp by ref code : %+v", err)
		return nil
	}

	result.OtpValue = newOtp
	result.ExpiredAt = time.Now().Add(time.Minute * 1)

	err = o.OtpRepository.OtpUpdate(tx, result)
	if err != nil {
		log.Printf("Failed update otp : %+v", err)
		return err
	}

	var toEmail string

	if result.UserRID.Valid {
		toEmail = result.UserRegister.Email

		subject, content, toEmail := mail.GetSenderParamEmailRegist(toEmail, newOtp)
		err := o.EmailSender.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
		if err != nil {
			log.Printf("Failed send email : %+v", err)
			return err
		}

		err = tx.Commit().Error
		if err != nil {
			return nil
		}

		return nil
	} else if result.UserID.Valid {
		toEmail = result.User.Email

		subject, content, toEmail := mail.GetSenderParamEmailRegist(toEmail, newOtp)
		err := o.EmailSender.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
		if err != nil {
			log.Printf("Failed send email : %+v", err)
			return err
		}

		err = tx.Commit().Error
		if err != nil {
			return nil
		}

		return nil
	}

	return util.EmailNotFound
}

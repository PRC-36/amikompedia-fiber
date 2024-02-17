package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/mail"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type RegisterUsecase interface {
	Register(ctx context.Context, request *request.UserRegisterRequest) (*response.RegisterResponseWithRefCode, error, []util.ApiErrorValidator)
}

type registerUsecaseImpl struct {
	DB                 *gorm.DB
	Validate           *validator.Validate
	EmailSender        mail.EmailSender
	RegisterRepository repository.RegisterRepository
	OtpRepository      repository.OtpRepository
}

func NewRegisterUsecase(db *gorm.DB, validate *validator.Validate, EmailSender mail.EmailSender, registerRepository repository.RegisterRepository, otpRepository repository.OtpRepository) RegisterUsecase {
	return &registerUsecaseImpl{DB: db, Validate: validate, EmailSender: EmailSender, RegisterRepository: registerRepository, OtpRepository: otpRepository}
}

func (r *registerUsecaseImpl) Register(ctx context.Context, requestData *request.UserRegisterRequest) (*response.RegisterResponseWithRefCode, error, []util.ApiErrorValidator) {

	tx := r.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := r.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body 1 : %+v", err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			log.Printf("Invalid request body 2 : %+v", ve)
			return nil, nil, util.HandleErrorValidator(ve)
		}
		log.Printf("Invalid request body 3 : %+v", err)
		return nil, err, nil
	}

	verified, _, _ := r.RegisterRepository.RegisterCheckEmailIsVerified(tx, requestData.Email)

	if verified {
		log.Printf("Email already used")
		return nil, util.EmailAlreadyUsed, nil
	}

	hashedPassword, err := util.HashPassword(requestData.Password)

	if err != nil {
		log.Printf("Failed hash password : %+v", err)
		return nil, err, nil
	}

	requestData.Password = hashedPassword
	requestRegisterEntity := requestData.ToEntity()

	err = r.RegisterRepository.RegisterCreate(tx, requestRegisterEntity)
	if err != nil {
		log.Printf("Failed create user register : %+v", err)
		return nil, err, nil
	}

	createRequestOtp := request.OtpCreateRequest{
		UserRID:   sql.NullInt32{Int32: int32(requestRegisterEntity.ID), Valid: true},
		OtpValue:  strconv.FormatInt(util.RandomInt(100000, 999999), 10),
		RefCode:   util.RandomCombineIntAndString(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}

	requestOtpEntity := createRequestOtp.ToEntity()

	err = r.OtpRepository.OtpCreate(tx, requestOtpEntity)
	if err != nil {
		log.Printf("Failed create otp  : %+v", err)
		return nil, err, nil
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err, nil
	}

	go func() {
		subject, content, toEmail := mail.GetSenderParamEmailRegist(requestData.Email, createRequestOtp.OtpValue)
		err := r.EmailSender.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
		if err != nil {
			log.Printf("Failed send email : %+v", err)
		}
	}()

	return requestRegisterEntity.ToUserRegisterResponse(requestOtpEntity.RefCode), nil, nil
}

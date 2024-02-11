package usecase

import (
	"context"
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"time"
)

type SessionUsecase interface {
	RenewAccessToken(ctx context.Context, requestData *request.TokenRenewRequest) (*response.TokenRenewResponse, error)
}

type sessionUsecase struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	TokenMaker        token.Maker
	ViperConfig       util.Config
	SessionRepository repository.SessionRepository
}

func NewSessionUsecase(db *gorm.DB, validate *validator.Validate, tokenMaker token.Maker,
	viperConfig util.Config, sessionRepository repository.SessionRepository) SessionUsecase {
	return &sessionUsecase{
		DB:                db,
		Validate:          validate,
		TokenMaker:        tokenMaker,
		ViperConfig:       viperConfig,
		SessionRepository: sessionRepository,
	}
}

func (t *sessionUsecase) RenewAccessToken(ctx context.Context, requestData *request.TokenRenewRequest) (*response.TokenRenewResponse, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := t.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	refreshPayload, err := t.TokenMaker.VerifyToken(requestData.RefreshToken, t.ViperConfig.TokenRefreshSymetricKey)

	if err != nil {
		log.Printf("Invalid refresh token : %+v", err)
		return nil, err
	}

	session, err := t.SessionRepository.FindByID(tx, refreshPayload.ID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Session id not found : %+v", err)
			return nil, util.SessionNotFound
		}
		log.Printf("Failed find session id : %+v", err)
		return nil, err
	}

	if session.IsBlocked {
		log.Printf("Session id is blocked : %+v", err)
		return nil, util.SessionIsBlocked
	}

	if session.UserID != refreshPayload.UserID {
		log.Printf("Session id not match with user id : %+v", err)
		return nil, util.SessionNotMatchUser
	}

	if session.RefreshToken != requestData.RefreshToken {
		log.Printf("Refresh token not match with session : %+v", err)
		return nil, util.InvalidRefreshToken
	}

	if time.Now().After(session.ExpiredAt) {
		log.Printf("Refresh token is expired : %+v", err)
		return nil, util.SessionExpired
	}

	accessToken, accessPayload, err := t.TokenMaker.CreateToken(token.UserPayload{
		UserID:   refreshPayload.UserID,
		Username: refreshPayload.Username,
	}, t.ViperConfig.TokenAccessSymetricKey, t.ViperConfig.TokenAccessDuration)

	if err != nil {
		log.Printf("Failed create access token : %+v", err)
		return nil, err
	}

	return &response.TokenRenewResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt.Format("2006-01-02 15:04:05"),
	}, nil
}

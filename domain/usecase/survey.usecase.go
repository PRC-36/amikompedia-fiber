package usecase

import (
	"context"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type SurveyUsecase interface {
	Create(ctx context.Context, userID string, request *request.SurveyRequest) (*response.SurveyResponse, error)
}

type surveyUsecaseImpl struct {
	DB               *gorm.DB
	Validate         *validator.Validate
	SurveyRepository repository.SurveyRepository
}

func NewSurveyUsecase(db *gorm.DB, validate *validator.Validate, surveyRepository repository.SurveyRepository) SurveyUsecase {
	return &surveyUsecaseImpl{DB: db, Validate: validate, SurveyRepository: surveyRepository}
}

func (s *surveyUsecaseImpl) Create(ctx context.Context, userID string, requestData *request.SurveyRequest) (*response.SurveyResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	requestSurveyEntity := requestData.ToEntity(userID)

	err = s.SurveyRepository.Create(tx, requestSurveyEntity)

	if err != nil {
		log.Printf("Failed create survey : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return requestSurveyEntity.ToSurveyResponse(), nil
}

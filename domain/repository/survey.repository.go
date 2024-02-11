package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type SurveyRepository interface {
	SurveyCreate(tx *gorm.DB, value *entity.Survey) error
}

type surveyRepositoryImpl struct {
}

func NewSurveyRepository() SurveyRepository {
	return &surveyRepositoryImpl{}
}

func (s *surveyRepositoryImpl) SurveyCreate(tx *gorm.DB, value *entity.Survey) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create survey : %v", result.Error))
		return result.Error
	}

	return nil
}

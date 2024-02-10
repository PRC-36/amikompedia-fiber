package repository

import (
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
)

type SurveyRepository interface {
	Create(tx *gorm.DB, value *entity.Survey) error
}

type surveyRepositoryImpl struct {
}

func NewSurveyRepository() SurveyRepository {
	return &surveyRepositoryImpl{}
}

func (s *surveyRepositoryImpl) Create(tx *gorm.DB, value *entity.Survey) error {
	result := tx.Create(value)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

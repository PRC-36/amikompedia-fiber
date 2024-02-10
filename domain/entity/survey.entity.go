package entity

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"time"
)

type Survey struct {
	ID               int       `gorm:"column:id" `
	UserID           string    `gorm:"column:user_id" `
	KnowsAmikompedia string    `gorm:"column:knows_amikompedia" `
	ImpressionDesc   string    `gorm:"column:impression_description" `
	CreatedAt        time.Time `gorm:"column:created_at" `
}

func (e *Survey) TableName() string {
	return "init_surveys"
}

func (e *Survey) ToSurveyResponse() *response.SurveyResponse {
	return &response.SurveyResponse{
		ID:               e.ID,
		UserID:           e.UserID,
		KnowsAmikompedia: e.KnowsAmikompedia,
		ImpressionDesc:   e.ImpressionDesc,
		CreatedAt:        e.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

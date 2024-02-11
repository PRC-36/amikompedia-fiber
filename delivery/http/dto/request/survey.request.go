package request

import "github.com/PRC-36/amikompedia-fiber/domain/entity"

type SurveyRequest struct {
	KnowsAmikompedia string `json:"knows_amikompedia" validate:"required"`
	ImpressionDesc   string `json:"impression_description" validate:"required"`
}

func (s SurveyRequest) ToEntity(userID string) *entity.Survey {
	return &entity.Survey{
		UserID:           userID,
		KnowsAmikompedia: s.KnowsAmikompedia,
		ImpressionDesc:   s.ImpressionDesc,
	}
}

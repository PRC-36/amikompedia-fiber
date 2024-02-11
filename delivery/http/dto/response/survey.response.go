package response

import "time"

type SurveyResponse struct {
	ID               int       `json:"id"`
	UserID           string    `json:"user_id"`
	KnowsAmikompedia string    `json:"knows_amikompedia"`
	ImpressionDesc   string    `json:"impression_description"`
	CreatedAt        time.Time `json:"created_at"`
}

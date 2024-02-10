package request

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"time"
)

type OtpCreateRequest struct {
	UserRID   sql.NullInt32  `json:"user_rid"`
	UserID    sql.NullString `json:"user_id"`
	OtpValue  string         `json:"otp_value"`
	RefCode   string         `json:"ref_code"`
	ExpiredAt time.Time      `json:"expired_at"`
}

func (r *OtpCreateRequest) ToEntity() *entity.Otp {
	return &entity.Otp{
		UserRID:   r.UserRID,
		UserID:    r.UserID,
		OtpValue:  r.OtpValue,
		RefCode:   r.RefCode,
		ExpiredAt: r.ExpiredAt,
	}
}

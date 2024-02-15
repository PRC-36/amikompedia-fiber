package request

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"time"
)

type OtpCreateRequest struct {
	UserRID   sql.NullInt32  `json:"user_rid" validate:"required"`
	UserID    sql.NullString `json:"user_id" validate:"required"`
	OtpValue  string         `json:"otp_value" validate:"required"`
	RefCode   string         `json:"ref_code" validate:"required"`
	ExpiredAt time.Time      `json:"expired_at" validate:"required"`
}

type OtpValidateRequest struct {
	OtpValue string `json:"otp_value" validate:"required,min=6,max=6"`
	RefCode  string `json:"ref_code" validate:"required"`
}

type OtpSendRequest struct {
	RefCode string `json:"ref_code" validate:"required"`
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

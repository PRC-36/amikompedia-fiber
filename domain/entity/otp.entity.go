package entity

import (
	"database/sql"
	"time"
)

type Otp struct {
	ID        int            `gorm:"primaryKey;autoIncrement;not null"`
	UserRID   sql.NullInt32  `gorm:"column:user_rid;not null"`
	UserID    sql.NullString `gorm:"column:user_id;not null"`
	OtpValue  string         `gorm:"column:otp_value;type:varchar(6);not null"`
	IsUsed    bool           `gorm:"column:is_used;type:boolean;not null;default:false"`
	RefCode   string         `gorm:"column:ref_code;not null"`
	ExpiredAt time.Time      `gorm:"column:expired_at;type:timestamp;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null"`
}

func (e *Otp) TableName() string {
	return "otps"
}

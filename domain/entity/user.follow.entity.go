package entity

import "time"

type UserFollow struct {
	ID          int       `gorm:"primaryKey"`
	FollowerID  string    `gorm:"column:follower_id"`
	FollowingID string    `gorm:"column:following_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (u *UserFollow) TableName() string {
	return "user_follows"
}

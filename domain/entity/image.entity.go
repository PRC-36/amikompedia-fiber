package entity

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"time"
)

type Image struct {
	ID        int            `gorm:"column:id"`
	UserID    string         `gorm:"column:user_uuid"`
	PostID    sql.NullString `gorm:"column:post_id"`
	ImageType string         `gorm:"column:image_type"`
	ImageUrl  string         `gorm:"column:image_url"`
	FilePath  string         `gorm:"column:file_path"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

func (i *Image) TableName() string {
	return "images"
}

func (i *Image) ToImageResponse() *response.ImageResponse {
	return &response.ImageResponse{
		ID:        i.ID,
		UserID:    i.UserID,
		PostID:    i.PostID.String,
		ImageType: i.ImageType,
		ImageUrl:  i.ImageUrl,
		FilePath:  i.FilePath,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func ToImageResponses(images []Image) []response.ImageResponse {
	var imageResponses []response.ImageResponse
	for _, image := range images {
		imageResponses = append(imageResponses, *image.ToImageResponse())
	}
	return imageResponses
}

package entity

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"log"
	"time"
)

type Post struct {
	ID        sql.NullString `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" `
	UserID    string         `gorm:"column:user_id"`
	RefPostID sql.NullString `gorm:"column:ref_post_id"`
	Content   string         `gorm:"column:content"`
	Images    []Image        `gorm:"foreignKey:PostID;references:ID"`
	User      User           `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
}

func (e *Post) TableName() string {
	return "posts"
}

func (e *Post) ToPostResponse() *response.PostResponse {
	return &response.PostResponse{
		ID:        e.ID.String,
		Content:   e.Content,
		User:      e.User.ToUserResponse(),
		Images:    ToImageResponses(e.Images),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func ToPostResponses(posts []Post, pagingMetadata *response.PostPageMetaData) *response.PostResponses {
	var postResponses []response.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, *post.ToPostResponse())
	}

	log.Printf("cek post", posts)
	return &response.PostResponses{
		Posts:  postResponses,
		Paging: pagingMetadata,
	}
}

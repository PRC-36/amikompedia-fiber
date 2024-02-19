package request

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
)

type PostRequest struct {
	Content string `json:"content" validate:"required"`
}

func (r *PostRequest) ToEntity(userId string) *entity.Post {
	return &entity.Post{
		UserID:  userId,
		Content: r.Content,
	}
}

type PostCommentRequest struct {
	Content   string `json:"content" validate:"required"`
	RefPostID string `json:"ref_post_id" validate:"required"`
}

func (r *PostCommentRequest) ToEntity(userId string) *entity.Post {
	return &entity.Post{
		UserID:    userId,
		Content:   r.Content,
		RefPostID: sql.NullString{Valid: true, String: r.RefPostID},
	}
}

type SearchPostRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page" validate:"min=1"`
	Size    int    `json:"size" validate:"min=1,max=100"`
}

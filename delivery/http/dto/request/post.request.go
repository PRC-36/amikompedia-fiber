package request

import (
	"database/sql"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
)

type PostRequest struct {
	Content   string         `json:"content" validate:"required"`
	UserID    string         `json:"-" validate:"required"`
	RefPostID sql.NullString `json:"ref_post_id"`
}

func (r *PostRequest) ToEntity() *entity.Post {
	return &entity.Post{
		UserID:    r.UserID,
		Content:   r.Content,
		RefPostID: r.RefPostID,
	}
}

type SearchPostRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page" validate:"min=1"`
	Size    int    `json:"size" validate:"min=1,max=100"`
}

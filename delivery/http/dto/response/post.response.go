package response

import (
	"time"
)

type PostResponse struct {
	ID        string          `json:"id"`
	RefPostID string          `json:"ref_post_id"`
	Content   string          `json:"content"`
	User      *UserResponse   `json:"user"`
	Images    []ImageResponse `json:"images"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type PostPageMetaData struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

type PostResponses struct {
	Posts  []PostResponse    `json:"posts"`
	Paging *PostPageMetaData `json:"paging"`
}

type DetailPostCommentResponse struct {
	Post     *PostResponse     `json:"post"`
	Comments []PostResponse    `json:"comments"`
	Paging   *PostPageMetaData `json:"paging"`
}

//func ToPostResponses(posts []entity.Post, pagingMetadata *PostPageMetaData) *PostResponses {
//	var postResponses []PostResponse
//	for _, post := range posts {
//		postResponses = append(postResponses, *post.ToPostResponse())
//	}
//
//	return &PostResponses{
//		Posts:  postResponses,
//		Paging: pagingMetadata,
//	}
//}

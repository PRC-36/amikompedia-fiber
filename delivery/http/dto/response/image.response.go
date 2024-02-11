package response

type ImageResponse struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	ImageType string `json:"image_type"`
	ImageUrl  string `json:"image_url"`
	FilePath  string `json:"file_path"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

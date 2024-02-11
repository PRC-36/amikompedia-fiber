package response

type UserResponse struct {
	ID        string          `json:"-"`
	Email     string          `json:"email"`
	Nim       string          `json:"nim"`
	Name      string          `json:"name"`
	Username  string          `json:"username"`
	Bio       string          `json:"bio"`
	Image     []ImageResponse `json:"images"`
	CreatedAt string          `json:"-"`
	UpdatedAt string          `json:"-"`
}

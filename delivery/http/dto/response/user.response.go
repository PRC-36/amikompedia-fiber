package response

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Nim       string `json:"nim"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

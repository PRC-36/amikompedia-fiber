package response

type RegisterResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Nim        string `json:"nim"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
}

type RegisterResponseWithRefCode struct {
	RefCode      string           `json:"ref_code"`
	RegisterUser RegisterResponse `json:"user_register"`
}

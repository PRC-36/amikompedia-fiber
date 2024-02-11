package response

type LoginResponse struct {
	Token *SessionsResponse `json:"token"`
	User  *UserResponse     `json:"user"`
}

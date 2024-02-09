package response

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UserLoginResponse struct {
	Token string       `json:"access_token"`
	User  UserResponse `json:"user"`
}

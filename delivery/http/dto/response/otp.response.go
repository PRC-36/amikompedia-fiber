package response

type OtpResponse struct {
	ID        int    `json:"id"`
	UserRID   int    `json:"user_rid"`
	UserID    string `json:"user_id"`
	ExpiredAt string `json:"expired_at"`
	CreatedAt string `json:"created_at"`
	RefCode   string `json:"ref_code"`
}

type OtpResponseWithToken struct {
	AccessToken string `json:"access_token"`
}

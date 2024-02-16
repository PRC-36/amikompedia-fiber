package response

type OtpResponse struct {
	ID        int    `json:"-"`
	UserRID   int    `json:"-"`
	UserID    string `json:"-"`
	ExpiredAt string `json:"-"`
	CreatedAt string `json:"-"`
	RefCode   string `json:"ref_code"`
}

type OtpResponseWithToken struct {
	AccessToken string `json:"access_token"`
}

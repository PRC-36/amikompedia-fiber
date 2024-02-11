package response

type TokenRenewResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
}

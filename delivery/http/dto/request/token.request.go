package request

type TokenRenewRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

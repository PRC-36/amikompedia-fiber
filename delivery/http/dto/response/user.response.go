package response

import "time"

type UserResponse struct {
	ID        string          `json:"-"`
	Email     string          `json:"email"`
	Nim       string          `json:"nim"`
	Name      string          `json:"name"`
	Username  string          `json:"username"`
	Bio       string          `json:"bio"`
	Image     []ImageResponse `json:"images"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
}

type ForgotPasswordResponseWithRefCode struct {
	RefCode string `json:"ref_code"`
}

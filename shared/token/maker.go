package token

type Maker interface {
	CreateToken(userPayload UserPayload) (string, error)
	VerifyToken(token string) (*Payload, error)
}

package token

import "time"

type Maker interface {
	CreateToken(userPayload UserPayload, secretKey string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string, secretKey string) (*Payload, error)
}

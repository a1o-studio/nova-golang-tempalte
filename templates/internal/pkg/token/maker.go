package token

import (
	"time"
)

const (
	SessionKey     = "shg_sid"
	AuthPayloadKey = "authorization_payload"
)

type Maker interface {
	CreateToken(userID int64, isStaff int16, duration time.Duration, tokenType TokenType) (string, *Payload, error)

	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}

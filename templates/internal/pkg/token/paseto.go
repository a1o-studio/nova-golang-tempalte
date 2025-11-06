package token

import (
	"fmt"
	"time"

	"github.com/a1ostudio/nova/internal/pkg/resp"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker 创建一个新的 PasetoMaker 实例
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userID int64, isStaff int16, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(userID, isStaff, duration, tokenType)
	if err != nil {
		return "", payload, err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, resp.ErrTokenInvalid
	}

	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

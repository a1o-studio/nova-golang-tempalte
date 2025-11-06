package token

import (
	"fmt"
	"time"

	"github.com/a1ostudio/nova/internal/pkg/resp"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType byte

const (
	TokenTypeAccess TokenType = iota + 1
	TokenTypeRefresh
)

type Payload struct {
	jwt.RegisteredClaims
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      TokenType `json:"token_type"` // TokenType 指示令牌的类型
	IsStaff   int16     `json:"is_staff"`   // is staff 0: normal user, 1: staff user
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userID int64, isStaff int16, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Type:      tokenType,
		IsStaff:   isStaff,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return payload, nil
}

// 检测 Payload 是否有效
func (payload *Payload) Valid(tokenType TokenType) error {
	if payload.Type != tokenType {
		return resp.ErrTokenInvalid
	}

	if time.Now().After(payload.ExpiredAt) {
		return resp.ErrTokenExpired
	}
	return nil
}

func (payload *Payload) GetIsStaff() bool {
	return payload.IsStaff == 1
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.ExpiredAt,
	}, nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: payload.IssuedAt,
	}, nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return "", nil
}

func (payload *Payload) GetSubject() (string, error) {
	return "", nil
}

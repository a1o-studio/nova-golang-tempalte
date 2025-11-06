package token

import (
	"testing"
	"time"

	"github.com/a1ostudio/nova/internal/pkg/resp"
	"github.com/a1ostudio/nova/internal/pkg/util"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	userID := util.RandomInt(1, 1000)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userID, 1, duration, TokenTypeAccess)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeAccess)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomInt(1, 1000), 1, -time.Minute, TokenTypeAccess)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeAccess)
	require.Error(t, err)
	require.EqualError(t, err, resp.ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestPasetoWrongTokenType(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomInt(1, 1000), 1, time.Minute, TokenTypeAccess)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeRefresh)
	require.Error(t, err)
	require.EqualError(t, err, resp.ErrTokenInvalid.Error())
	require.Nil(t, payload)
}

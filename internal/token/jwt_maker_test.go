package token

import (
	"github.com/kavinddd/mangtoon_be/internal/role"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwtMaker(t *testing.T) {
	tokenMaker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomUsername()
	roles := []string{role.Reader}
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	jwt, payload, err := tokenMaker.CreateToken(username, roles, duration)

	require.NoError(t, err)
	require.NotEmpty(t, jwt)
	require.NotEmpty(t, payload)

	payload, err = tokenMaker.VerifyToken(jwt)

	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJwt(t *testing.T) {
	tokenMaker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	jwt, payload, err := tokenMaker.CreateToken(util.RandomUsername(), []string{role.Reader}, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwt)
	require.NotEmpty(t, payload)

	payload, err = tokenMaker.VerifyToken(jwt)
	require.Error(t, err)
	require.EqualError(t, err, "token has invalid claims: token is expired")
	require.Nil(t, payload)

}

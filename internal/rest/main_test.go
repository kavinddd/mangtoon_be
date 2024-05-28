package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kavinddd/mangtoon_be/internal/db"
	"github.com/kavinddd/mangtoon_be/internal/token"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func addAuthorizationToHeader(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
	roles []string,
) {
	token, payload, err := tokenMaker.CreateToken(username, roles, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

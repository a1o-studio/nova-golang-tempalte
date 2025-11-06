package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/a1ostudio/nova/internal/config"
	"github.com/a1ostudio/nova/internal/pkg/token"
	"github.com/a1ostudio/nova/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type server struct {
	config     config.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func newTestServer(t *testing.T) *server {
	config := config.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		LimitRate:           1,
		LimitBurst:          1,
	}

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	return &server{
		config:     config,
		tokenMaker: tokenMaker,
		router:     gin.New(),
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

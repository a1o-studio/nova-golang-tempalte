package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestTimeout(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, response *httptest.ResponseRecorder, code int)
	}{
		{
			name: "Timeout",
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder, code int) {
				require.Equal(t, code, response.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			server.router.Use(Timeout(50 * time.Millisecond))
			server.router.GET("/slow", func(c *gin.Context) {
				time.Sleep(100 * time.Millisecond) // 模拟一个慢请求
				c.String(http.StatusOK, "ok")
			})
			server.router.GET("/fast", func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			})

			recorder := httptest.NewRecorder()
			slow, err := http.NewRequest(http.MethodGet, "/slow", nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, slow)
			tc.checkResponse(t, recorder, http.StatusGatewayTimeout)

			fast, err := http.NewRequest(http.MethodGet, "/fast", nil)
			require.NoError(t, err)
			recorder = httptest.NewRecorder()
			server.router.ServeHTTP(recorder, fast)
			tc.checkResponse(t, recorder, http.StatusOK)
		})
	}
}

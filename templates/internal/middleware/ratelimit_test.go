package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRateLimitByIPMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, response *httptest.ResponseRecorder, statusCode int)
	}{
		{
			name: "Limit",
			checkResponse: func(t *testing.T, response *httptest.ResponseRecorder, statusCode int) {
				require.Equal(t, statusCode, response.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			path := "/test"
			server.router.GET(
				path,
				RateLimitByIPMiddleware(server.config.LimitRate, server.config.LimitBurst),
				func(c *gin.Context) {
					c.String(200, "ok")
				},
			)

			request, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)
			request.RemoteAddr = "127.0.0.1:12345"

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, http.StatusOK)

			// 模拟多次请求以触发限流
			for j := 0; j < server.config.LimitBurst+1; j++ {
				recorder = httptest.NewRecorder()
				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder, http.StatusTooManyRequests)
			}

			// 等待 1s 之后重试
			time.Sleep(time.Second)
			recorder = httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, http.StatusOK)
		})
	}
}

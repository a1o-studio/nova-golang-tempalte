package middleware

import (
	"context"
	"time"

	"github.com/a1ostudio/nova/internal/pkg/resp"

	"github.com/gin-gonic/gin"
)

func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				resp.TimeoutError(c, resp.WithMessage("request timed out"))
				c.Abort()
			}
		case <-done:
			// 请求完成
		}
	}
}

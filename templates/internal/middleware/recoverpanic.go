package middleware

import (
	"fmt"

	"github.com/a1ostudio/nova/internal/pkg/resp"

	"github.com/gin-gonic/gin"
)

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				resp.ServerError(c, resp.WithMessage(msg))
			}
			// 终端后续中间件
			c.Abort()
		}()

		c.Next()
	}
}

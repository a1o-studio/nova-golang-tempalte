package middleware

import (
	"sync"
	"time"

	"github.com/a1ostudio/nova/internal/pkg/resp"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var clients sync.Map

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitByIPMiddleware returns a middleware that limits the number of requests
// from a single IP address to 'r' requests per 'b' seconds.
func RateLimitByIPMiddleware(r, b int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		limiter, _ := clients.LoadOrStore(ip, &client{
			limiter:  rate.NewLimiter(rate.Limit(r), b),
			lastSeen: time.Now(),
		})
		rl := limiter.(*client)

		if !rl.limiter.Allow() {
			resp.TooManyRequestsError(ctx)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CleanupClients(interval time.Duration, expiration time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		clients.Range(func(key, value any) bool {
			rl := value.(*client)
			if now.Sub(rl.lastSeen) > expiration {
				clients.Delete(key)
			}
			return true
		})
	}
}

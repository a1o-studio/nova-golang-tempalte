package logger

import (
	"sync"
	"time"

	"github.com/a1ostudio/nova/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	log       *zap.Logger
	atomicLvl zap.AtomicLevel
	once      sync.Once
)

// Init initializes the global logger.
func NewLogger(env config.Env) {
	once.Do(func() {
		level := zap.InfoLevel
		if env == config.Dev {
			level = zap.DebugLevel
		}
		atomicLvl = zap.NewAtomicLevelAt(level)

		cfg := zap.Config{
			Level:            atomicLvl,
			Development:      env == config.Dev,
			Encoding:         "json",
			OutputPaths:      []string{"stdout", "logs/app.log"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:      "time",
				LevelKey:     "level",
				MessageKey:   "message",
				CallerKey:    "caller",
				EncodeLevel:  zapcore.CapitalLevelEncoder,
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}

		if env == "dev" {
			cfg.Encoding = "console"
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			cfg.OutputPaths = []string{"stdout"}
		}

		logger, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		log = logger
		zap.ReplaceGlobals(log)
	})
}

// L returns the global logger.
func L() *zap.Logger {
	return log
}

// WithRequest returns a logger with request_id field.
func WithRequest(requestID string) *zap.Logger {
	return log.With(zap.String("request_id", requestID))
}

// WithFields adds custom fields to the logger.
func WithFields(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// SetLevel allows dynamic log level adjustment.
func SetLevel(level zapcore.Level) {
	atomicLvl.SetLevel(level)
}

// Sync flushes the logger's buffer.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// shortPath trims query strings and long IDs in URL paths for cleaner logging.
func shortPath(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		return path
	}
	parts := strings.Split(u.Path, "/")
	for i, part := range parts {
		if len(part) > 24 {
			parts[i] = "...id"
		}
	}
	return strings.Join(parts, "/")
}

// emojiStatus returns an emoji based on HTTP status code.
func emojiStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "✅"
	case code >= 400 && code < 500:
		return "⚠️"
	case code >= 500:
		return "❌"
	default:
		return "🔍"
	}
}

// LoggerMiddleware logs simplified request info with emoji and path trimming.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		logger := WithRequest(requestID)
		c.Set("logger", logger)
		c.Writer.Header().Set("X-Request-ID", requestID)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		method := c.Request.Method
		path := shortPath(c.Request.URL.Path)
		status := c.Writer.Status()
		emoji := emojiStatus(status)

		logger.Info(emoji+" HTTP",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
			zap.String("duration", duration.String()),
		)
	}
}

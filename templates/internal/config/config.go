package config

import (
	"time"

	"github.com/spf13/viper"
)

type Env = string //	@name	Env

const (
	Dev     Env = "dev"
	Staging Env = "staging"
	Prod    Env = "prod"
)

type Config struct {
	Domain               string        `mapstructure:"SERVER_DOMAIN"` // 服务器地址
	Port                 int           `mapstructure:"SERVER_PORT"`
	Env                  Env           `mapstructure:"ENV"` // dev, staging, prod, etc.
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`  // 数据库迁移地址
	RedisPassword        string        `mapstructure:"REDIS_PASSWORD"` // Redis 密码
	RedisPort            int           `mapstructure:"REDIS_PORT"`     // Redis 地址
	LimitRate            int           `mapstructure:"LIMIT_RATE"`     // 每秒允许的请求数
	LimitBurst           int           `mapstructure:"LIMIT_BURST"`    // 允许的突发请求数
	AesSecret            string        `mapstructure:"AES_SECRET"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`    // Token 对称密钥
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`  // 访问令牌有效期
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"` // 刷新令牌有效期
	InvitationDuration   time.Duration `mapstructure:"INVITATION_DURATION"`    // 邀请函有效期
	ImportsPath          string        `mapstructure:"IMPORTS_PATH"`

	// 分布式锁配置参数
	LockTTL         time.Duration `mapstructure:"LOCK_TTL"`          // 锁的生存时间，默认 2s
	MaxWaitTime     time.Duration `mapstructure:"MAX_WAIT_TIME"`     // 等待锁的最大时间，默认 1s
	InitialWaitTime time.Duration `mapstructure:"INITIAL_WAIT_TIME"` // 初始等待时间，默认 50ms
	MaxSingleWait   time.Duration `mapstructure:"MAX_SINGLE_WAIT"`   // 单次最大等待时间，默认 200ms
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// 设置分布式锁参数默认值
	viper.SetDefault("LOCK_TTL", "2s")
	viper.SetDefault("MAX_WAIT_TIME", "1s")
	viper.SetDefault("INITIAL_WAIT_TIME", "50ms")
	viper.SetDefault("MAX_SINGLE_WAIT", "200ms")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

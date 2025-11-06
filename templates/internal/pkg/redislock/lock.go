package redislock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// AcquireLock tries to create a lock with a random owner value. It returns
// (owner, true, nil) when the lock is acquired. If not acquired returns
// ("", false, nil). Any redis error is returned as err.
func AcquireLock(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (string, bool, error) {
	owner := uuid.NewString()
	ok, err := rdb.SetNX(ctx, key, owner, ttl).Result()
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return owner, true, nil
}

// ReleaseLock releases the lock only if the owner matches. It uses a Lua
// script to ensure atomicity. Returns true if the lock was released.
func ReleaseLock(ctx context.Context, rdb *redis.Client, key string, owner string) (bool, error) {
	// Lua script: if value matches then del and return 1 else return 0
	const lua = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`
	res, err := rdb.Eval(ctx, lua, []string{key}, owner).Result()
	if err != nil {
		return false, err
	}
	if v, ok := res.(int64); ok && v > 0 {
		return true, nil
	}
	return false, nil
}

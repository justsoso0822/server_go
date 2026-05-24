package lock

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"server_go/internal/autodb"
	"server_go/utility/dbcache"

	"github.com/gogf/gf/v2/database/gredis"
)

const (
	acquireTimeoutMs = 1000
	retryBaseMs      = 20
	retryMaxMs       = 200
	lockTTLMs        = 30000
)

func Lock(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("[Lock] key is required")
	}
	redis := autodb.Redis(ctx)
	redisKey := dbcache.BuildKey(ctx, "lock", key)
	token := fmt.Sprintf("%d:%d:%d", os.Getpid(), time.Now().UnixNano(), rand.Int63())

	deadline := time.Now().Add(time.Duration(acquireTimeoutMs) * time.Millisecond)
	retryCount := 0

	for time.Now().Before(deadline) {
		ok, err := tryAcquire(ctx, redis, redisKey, token)
		if err != nil {
			return "", err
		}
		if ok {
			return token, nil
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		cap := math.Min(float64(retryMaxMs), float64(retryBaseMs)*math.Pow(2, float64(retryCount)))
		sleepMs := int(math.Min(float64(remaining.Milliseconds()), float64(rand.Intn(int(cap)+1))))
		retryCount++
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	return "", nil
}

func Unlock(ctx context.Context, key, token string) error {
	if key == "" || token == "" {
		return nil
	}
	redis := autodb.Redis(ctx)
	redisKey := dbcache.BuildKey(ctx, "lock", key)
	script := `if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end`
	_, err := redis.Do(ctx, "EVAL", script, 1, redisKey, token)
	return err
}

func tryAcquire(ctx context.Context, redis *gredis.Redis, key, token string) (bool, error) {
	res, err := redis.Do(ctx, "SET", key, token, "PX", lockTTLMs, "NX")
	if err != nil {
		return false, err
	}
	return res.String() == "OK", nil
}

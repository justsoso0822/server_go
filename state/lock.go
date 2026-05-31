package state

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"server_go/tools/autodb"
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
	rc := autodb.Redis(ctx)
	redisKey := BuildKey(ctx, "lock", key)
	token := fmt.Sprintf("%d:%d:%d", os.Getpid(), time.Now().UnixNano(), rand.Int63())

	deadline := time.Now().Add(time.Duration(acquireTimeoutMs) * time.Millisecond)
	retryCount := 0

	for time.Now().Before(deadline) {
		ok, err := rc.SetNX(ctx, redisKey, token, time.Duration(lockTTLMs)*time.Millisecond).Result()
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
	rc := autodb.Redis(ctx)
	redisKey := BuildKey(ctx, "lock", key)
	script := `if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end`
	return rc.Eval(ctx, script, []string{redisKey}, token).Err()
}

func BuildKey(ctx context.Context, parts ...string) string {
	channel := autodb.GetChannel(ctx)
	if channel == "" {
		channel = autodb.DefaultChannelName
	}
	result := channel
	for _, p := range parts {
		result += ":" + p
	}
	return result
}

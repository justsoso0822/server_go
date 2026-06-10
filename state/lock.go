package state

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"server_go/tools/autodb"
)

const (
	// 最多等待 1 秒拿锁，超过就让业务返回“系统繁忙”，避免请求线程长时间堆积。
	acquireTimeoutMs = 1000
	// 重试间隔从 20ms 开始指数增长，并设置上限，兼顾快速成功和热点保护。
	retryBaseMs = 20
	retryMaxMs  = 200
	// 锁租约兜底 30 秒；即使持锁进程崩溃，Redis 也会自动释放。
	lockTTLMs = 30000
)

var ErrLockNotAcquired = errors.New("系统繁忙，请稍后再试")

// 使用 Redis 做短临界区互斥，成功时返回释放锁必须携带的 token。
// 没有拿到锁时显式返回 ErrLockNotAcquired 和空 token，避免出现 token == "" 且 err == nil 的模糊状态。
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
	return "", ErrLockNotAcquired
}

// 释放锁时校验 token，确保只释放自己持有的那把锁。
func Unlock(ctx context.Context, key, token string) error {
	if key == "" || token == "" {
		return nil
	}
	rc := autodb.Redis(ctx)
	redisKey := BuildKey(ctx, "lock", key)
	script := `if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end`
	return rc.Eval(ctx, script, []string{redisKey}, token).Err()
}

// 拼接项目内统一的 Redis key，并自动补上渠道前缀。
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

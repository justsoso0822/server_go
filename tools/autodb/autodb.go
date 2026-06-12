package autodb

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"server_go/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type channelContextKey struct{}
type requestIDContextKey struct{}
type logUIDContextKey struct{}
type logOpenIDContextKey struct{}

var channelKey = channelContextKey{}
var reqIDKey = requestIDContextKey{}
var logUIDKey = logUIDContextKey{}
var logOpenIDKey = logOpenIDContextKey{}

const DefaultChannelName = "default"
const connectTimeout = 5 * time.Second

var reservedChannelNames = map[string]struct{}{
	"health":   {},
	"internal": {},
	"api":      {},
	"other":    {},
	"test":     {},
}

var (
	mu                 sync.RWMutex
	configuredChannels = map[string]struct{}{}
	channelList        []string
	dbs                = map[string]*gorm.DB{}
	redisClients       = map[string]*redis.Client{}
	cacheEnabled       = map[string]bool{}
	sf                 singleflight.Group
)

func Init(cfg *config.Config, log *zap.Logger) error {
	// go-redis 内部连接池告警默认写标准库 log。接到 zap 后，线上日志能保持 JSON/字段化。
	redis.SetLogger(newRedisLogger(log))

	if _, ok := cfg.Database[DefaultChannelName]; !ok {
		return fmt.Errorf("default channel (database.default + redis.default) is required")
	}

	set := make(map[string]struct{}, len(cfg.Database))
	list := make([]string, 0, len(cfg.Database))

	for name := range cfg.Database {
		if _, hit := reservedChannelNames[name]; hit {
			return fmt.Errorf("channel name %q is reserved by routing layer", name)
		}
		if _, ok := cfg.Redis[name]; !ok {
			return fmt.Errorf("channel %q requires redis.%s config", name, name)
		}
		set[name] = struct{}{}
		list = append(list, name)
	}

	sort.Strings(list)

	nextDBs := make(map[string]*gorm.DB, len(list))
	nextRedisClients := make(map[string]*redis.Client, len(list))
	nextCacheEnabled := make(map[string]bool, len(list))
	cleanup := func() {
		for _, db := range nextDBs {
			if sqlDB, err := db.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}
		for _, rc := range nextRedisClients {
			_ = rc.Close()
		}
	}

	for _, name := range list {
		dbCfg := cfg.Database[name]
		db, err := openDB(name, dbCfg, log)
		if err != nil {
			cleanup()
			return fmt.Errorf("open db %s: %w", name, err)
		}
		nextDBs[name] = db

		rc, err := openRedis(name, cfg.Redis[name])
		if err != nil {
			cleanup()
			return err
		}
		nextRedisClients[name] = rc
		nextCacheEnabled[name] = dbCfg.Cache
	}

	mu.Lock()
	dbs = nextDBs
	redisClients = nextRedisClients
	cacheEnabled = nextCacheEnabled
	configuredChannels = set
	channelList = list
	mu.Unlock()

	return nil
}

func openDB(name string, cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := strings.TrimPrefix(cfg.Link, "mysql:")

	// gorm.Open 返回的是带连接池管理能力的 *gorm.DB。真正的池参数在下面通过 db.DB()
	// 拿到底层 *sql.DB 后配置；GORM 本身不负责 MaxOpen/MaxIdle 这些细节。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 自定义 logger 实现 GORM logger.Interface，可以把慢查询、错误 SQL 与 request_id 关联起来。
		Logger: newGormLogger(log, name, cfg.Debug, cfg.SlowThreshold),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	if err := sqlDB.PingContext(pingCtx); err != nil {
		cancel()
		_ = sqlDB.Close()
		return nil, err
	}
	cancel()

	maxIdle := cfg.MaxIdle
	if maxIdle <= 0 {
		maxIdle = 10
	}
	maxOpen := cfg.MaxOpen
	if maxOpen <= 0 {
		maxOpen = 100
	}
	maxLifetime := cfg.MaxLifetime
	if maxLifetime <= 0 {
		maxLifetime = 3600
	}

	// MaxIdle 控制空闲连接数，MaxOpen 控制总连接上限。
	// ConnMaxLifetime 应小于 MySQL/代理层的空闲回收时间，减少服务端主动断开造成的坏连接。
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	return db, nil
}

func openRedis(name string, cfg config.RedisConfig) (*redis.Client, error) {
	poolSize := cfg.PoolSize
	if poolSize <= 0 {
		poolSize = 100
	}
	minIdleConns := cfg.MinIdleConns
	if minIdleConns <= 0 {
		minIdleConns = 10
	}

	// go-redis Client 内部是并发安全的连接池，一个 channel 复用一个 Client 即可。
	// 超时和连接池参数显式配置，避免 Redis 抖动时请求无限排队或被慢连接拖住。
	rc := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Pass,
		DB:           cfg.DB,
		DialTimeout:  durationFromMillis(cfg.DialTimeoutMs, 3000),
		ReadTimeout:  durationFromMillis(cfg.ReadTimeoutMs, 500),
		WriteTimeout: durationFromMillis(cfg.WriteTimeoutMs, 500),
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		PoolTimeout:  durationFromMillis(cfg.PoolTimeoutMs, 1000),
	})
	pingCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if err := rc.Ping(pingCtx).Err(); err != nil {
		_ = rc.Close()
		return nil, fmt.Errorf("ping redis %s: %w", name, err)
	}
	return rc, nil
}

func durationFromMillis(value, fallback int) time.Duration {
	if value <= 0 {
		value = fallback
	}
	return time.Duration(value) * time.Millisecond
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	var firstErr error
	for name, db := range dbs {
		sqlDB, err := db.DB()
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("get sql db %s: %w", name, err)
			}
			continue
		}
		if err := sqlDB.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close db %s: %w", name, err)
		}
	}
	for name, rc := range redisClients {
		if err := rc.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close redis %s: %w", name, err)
		}
	}

	dbs = map[string]*gorm.DB{}
	redisClients = map[string]*redis.Client{}
	cacheEnabled = map[string]bool{}
	configuredChannels = map[string]struct{}{}
	channelList = nil

	return firstErr
}

func ConfiguredChannels() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, len(channelList))
	copy(out, channelList)
	return out
}

func IsConfiguredChannel(channel string) bool {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return false
	}
	mu.RLock()
	_, ok := configuredChannels[channel]
	mu.RUnlock()
	return ok
}

func WithChannel(ctx context.Context, channel string) context.Context {
	return context.WithValue(ctx, channelKey, strings.TrimSpace(channel))
}

func GetChannel(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(channelKey); v != nil {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func BackgroundWithChannel(ctx context.Context) context.Context {
	bgCtx := context.Background()
	if channel := GetChannel(ctx); channel != "" {
		bgCtx = WithChannel(bgCtx, channel)
	}
	if requestID := GetRequestID(ctx); requestID != "" {
		bgCtx = WithRequestID(bgCtx, requestID)
	}
	if uid := GetLogUID(ctx); uid != 0 {
		bgCtx = WithLogUID(bgCtx, uid)
	}
	if openid := GetLogOpenID(ctx); openid != "" {
		bgCtx = WithLogOpenID(bgCtx, openid)
	}
	return bgCtx
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, reqIDKey, strings.TrimSpace(requestID))
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(reqIDKey); v != nil {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func WithLogIdentity(ctx context.Context, uid int64, openid string) context.Context {
	if uid != 0 {
		ctx = WithLogUID(ctx, uid)
	}
	if openid != "" {
		ctx = WithLogOpenID(ctx, openid)
	}
	return ctx
}

func WithLogUID(ctx context.Context, uid int64) context.Context {
	return context.WithValue(ctx, logUIDKey, uid)
}

func GetLogUID(ctx context.Context) int64 {
	if ctx == nil {
		return 0
	}
	if v := ctx.Value(logUIDKey); v != nil {
		if uid, ok := v.(int64); ok {
			return uid
		}
	}
	return 0
}

func WithLogOpenID(ctx context.Context, openid string) context.Context {
	return context.WithValue(ctx, logOpenIDKey, strings.TrimSpace(openid))
}

func GetLogOpenID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(logOpenIDKey); v != nil {
		if openid, ok := v.(string); ok {
			return strings.TrimSpace(openid)
		}
	}
	return ""
}

func DB(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	db := dbs[channel]
	mu.RUnlock()
	if db == nil {
		return nil, fmt.Errorf("database channel %q is not initialized", channel)
	}
	// WithContext 不会新建连接，它只是把 context 绑定到后续 SQL：
	// 一方面支持超时/取消，另一方面 GORM logger.Trace 能从 context 取 request_id。
	return db.WithContext(ctx), nil
}

func Redis(ctx context.Context) (*redis.Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	rc := redisClients[channel]
	mu.RUnlock()
	if rc == nil {
		return nil, fmt.Errorf("redis channel %q is not initialized", channel)
	}
	return rc, nil
}

func CacheEnabled(ctx context.Context) bool {
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	enabled := cacheEnabled[channel]
	mu.RUnlock()
	return enabled
}

func BuildCacheKey(ctx context.Context, key string) string {
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	return "mysql_cache:" + channel + ":" + strings.TrimSpace(key)
}

func Cache[T any](ctx context.Context, key string, ttl time.Duration, load func() (T, error)) (T, error) {
	var zero T
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(key) == "" || !CacheEnabled(ctx) {
		return load()
	}
	rc, err := Redis(ctx)
	if err != nil {
		return load()
	}
	fullKey := BuildCacheKey(ctx, key)
	// 第一层直接读缓存，命中时不进入 singleflight，减少锁竞争。
	data, err := rc.Get(ctx, fullKey).Bytes()
	if err == nil {
		var dest T
		if json.Unmarshal(data, &dest) == nil {
			return dest, nil
		}
	}
	result, err, _ := sf.Do(fullKey, func() (any, error) {
		// singleflight 只合并同进程内相同 key 的并发回源；多实例之间仍可能同时回源。
		// 如果以后遇到跨实例缓存击穿，可以在这里再叠加 Redis 分布式锁。
		data, err := rc.Get(ctx, fullKey).Bytes()
		if err == nil {
			var dest T
			if json.Unmarshal(data, &dest) == nil {
				return dest, nil
			}
		}

		dest, err := load()
		if err != nil {
			return zero, err
		}
		if ttl > 0 {
			data, err := json.Marshal(dest)
			if err != nil {
				return dest, nil
			}
			_ = rc.Set(ctx, fullKey, data, ttl).Err()
		}
		return dest, nil
	})
	if err != nil {
		return zero, err
	}
	dest, ok := result.(T)
	if !ok {
		return zero, fmt.Errorf("cache result type mismatch for key %s", fullKey)
	}
	return dest, nil
}

func DelCache(ctx context.Context, keys ...string) {
	if len(keys) == 0 || !CacheEnabled(ctx) {
		return
	}
	rc, err := Redis(ctx)
	if err != nil {
		return
	}
	cacheKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		if strings.TrimSpace(key) != "" {
			cacheKeys = append(cacheKeys, BuildCacheKey(ctx, key))
		}
	}
	if len(cacheKeys) > 0 {
		_ = rc.Del(ctx, cacheKeys...).Err()
	}
}

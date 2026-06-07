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
	// link 格式: "mysql:user:pass@tcp(host:port)/dbname?params"
	dsn := strings.TrimPrefix(cfg.Link, "mysql:")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	return db, nil
}

func openRedis(name string, cfg config.RedisConfig) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Pass,
		DB:       cfg.DB,
	})
	pingCtx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if err := rc.Ping(pingCtx).Err(); err != nil {
		_ = rc.Close()
		return nil, fmt.Errorf("ping redis %s: %w", name, err)
	}
	return rc, nil
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

func DB(ctx context.Context) *gorm.DB {
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	db := dbs[channel]
	mu.RUnlock()
	if db == nil {
		return nil
	}
	return db.WithContext(ctx)
}

func Redis(ctx context.Context) *redis.Client {
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	rc := redisClients[channel]
	mu.RUnlock()
	return rc
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
	rc := Redis(ctx)
	if rc == nil {
		return load()
	}
	fullKey := BuildCacheKey(ctx, key)
	data, err := rc.Get(ctx, fullKey).Bytes()
	if err == nil {
		var dest T
		if json.Unmarshal(data, &dest) == nil {
			return dest, nil
		}
	}
	result, err, _ := sf.Do(fullKey, func() (any, error) {
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
	rc := Redis(ctx)
	if rc == nil {
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

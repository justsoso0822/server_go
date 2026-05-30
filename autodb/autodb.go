package autodb

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"server_gin/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type channelKey struct{}

var ctxKey = channelKey{}

const DefaultChannelName = "default"

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
)

func Init(cfg *config.Config) error {
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
		db, err := openDB(dbCfg)
		if err != nil {
			cleanup()
			return fmt.Errorf("open db %s: %w", name, err)
		}
		nextDBs[name] = db

		rCfg := cfg.Redis[name]
		rc := redis.NewClient(&redis.Options{
			Addr:     rCfg.Address,
			Password: rCfg.Pass,
			DB:       rCfg.DB,
		})
		nextRedisClients[name] = rc
	}

	mu.Lock()
	dbs = nextDBs
	redisClients = nextRedisClients
	configuredChannels = set
	channelList = list
	mu.Unlock()

	return nil
}

func openDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// link 格式: "mysql:user:pass@tcp(host:port)/dbname?params"
	dsn := strings.TrimPrefix(cfg.Link, "mysql:")

	logLevel := logger.Silent
	if cfg.Debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

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
	return context.WithValue(ctx, ctxKey, strings.TrimSpace(channel))
}

func GetChannel(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(ctxKey); v != nil {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func BackgroundWithChannel(ctx context.Context) context.Context {
	channel := GetChannel(ctx)
	if channel == "" {
		return context.Background()
	}
	return WithChannel(context.Background(), channel)
}

func DB(ctx context.Context) *gorm.DB {
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	db, ok := dbs[channel]
	mu.RUnlock()
	if !ok {
		return dbs[DefaultChannelName]
	}
	return db.WithContext(ctx)
}

func Redis(ctx context.Context) *redis.Client {
	channel := GetChannel(ctx)
	if channel == "" {
		channel = DefaultChannelName
	}
	mu.RLock()
	rc, ok := redisClients[channel]
	mu.RUnlock()
	if !ok {
		return redisClients[DefaultChannelName]
	}
	return rc
}

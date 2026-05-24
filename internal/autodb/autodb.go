package autodb

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type channelKey struct{}

var ctxKey = channelKey{}

// DefaultChannelName 是无渠道前缀请求默认绑定的 channel 名。
const DefaultChannelName = "default"

// reservedChannelNames 是 HTTP 路由层占用的保留前缀，不允许作为 channel 名使用，
// 否则会与 /health、/internal/control 等固定路由或渠道内部子路径冲突。
var reservedChannelNames = map[string]struct{}{
	"health":   {},
	"internal": {},
	"api":      {},
	"other":    {},
	"test":     {},
}

var (
	channelsMu         sync.RWMutex
	configuredChannels = map[string]struct{}{}
	channelList        []string
	channelsLoaded     bool
)

// LoadConfiguredChannels 启动期扫描 database.* 与 redis.*，构建合法 channel 集合。
// 要求：每个 database.<name> 必须有同名 redis.<name>；name 不在保留前缀列表中；必须包含 default。
// 仅在进程启动时调用一次，运行期不再读配置。
func LoadConfiguredChannels(ctx context.Context) ([]string, error) {
	dbVar, err := g.Cfg().Get(ctx, "database")
	if err != nil {
		return nil, fmt.Errorf("read database config: %w", err)
	}
	if dbVar.IsNil() {
		return nil, fmt.Errorf("database config is missing")
	}

	dbMap := dbVar.Map()
	set := make(map[string]struct{}, len(dbMap))
	list := make([]string, 0, len(dbMap))

	for raw := range dbMap {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, hit := reservedChannelNames[name]; hit {
			return nil, fmt.Errorf("channel name %q is reserved by routing layer", name)
		}
		redisVar, err := g.Cfg().Get(ctx, "redis."+name)
		if err != nil {
			return nil, fmt.Errorf("read redis.%s: %w", name, err)
		}
		if redisVar.IsNil() {
			return nil, fmt.Errorf("channel %q requires redis.%s config", name, name)
		}
		set[name] = struct{}{}
		list = append(list, name)
	}

	if _, ok := set[DefaultChannelName]; !ok {
		return nil, fmt.Errorf("default channel (database.%s + redis.%s) is required", DefaultChannelName, DefaultChannelName)
	}

	sort.Strings(list)

	channelsMu.Lock()
	configuredChannels = set
	channelList = list
	channelsLoaded = true
	channelsMu.Unlock()

	return list, nil
}

// ConfiguredChannels 返回启动期构建的合法 channel 列表（有序拷贝）。
func ConfiguredChannels() []string {
	channelsMu.RLock()
	defer channelsMu.RUnlock()
	out := make([]string, len(channelList))
	copy(out, channelList)
	return out
}

// WithChannel 将请求的 channel 写入上下文，供 DAO、Redis 和原生 SQL 统一读取。
func WithChannel(ctx context.Context, channel string) context.Context {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxKey, channel)
}

// GetChannel 优先从上下文取 channel，其次回退到路由参数。
func GetChannel(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(ctxKey); v != nil {
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s)
		}
	}
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		if s := strings.TrimSpace(r.GetRouter("channel").String()); s != "" {
			return s
		}
		if s := strings.TrimSpace(r.GetParam("channel").String()); s != "" {
			return s
		}
	}
	return ""
}

// BackgroundWithChannel 创建后台 context，保留当前请求的 channel 信息。
// 用于异步 goroutine 中保持渠道分流。
func BackgroundWithChannel(ctx context.Context) context.Context {
	channel := GetChannel(ctx)
	if channel == "" {
		return context.Background()
	}
	return WithChannel(context.Background(), channel)
}

// DB 返回数据库实例；未指定 group 或 group 为 default 时，优先使用当前请求的 channel。
func DB(ctx context.Context, groups ...string) gdb.DB {
	group := ""
	if len(groups) > 0 {
		group = groups[0]
	}
	return g.DB(resolveGroup(ctx, group, gdb.DefaultGroupName))
}

// DefaultDB 返回框架默认数据库实例，不读取请求 channel。
func DefaultDB() gdb.DB {
	return g.DB(gdb.DefaultGroupName)
}

// Model 返回已绑定请求上下文的安全 Model，并按 channel 选择数据库分组。
func Model(ctx context.Context, table string, groups ...string) *gdb.Model {
	return DB(ctx, groups...).Model(table).Safe().Ctx(ctx)
}

// Transaction 在当前 channel 对应的数据库分组中执行事务。
func Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error, groups ...string) error {
	return DB(ctx, groups...).Transaction(ctx, f)
}

// Redis 返回 Redis 实例；未指定 group 或 group 为 default 时，优先使用当前请求的 channel。
func Redis(ctx context.Context, groups ...string) *gredis.Redis {
	group := ""
	if len(groups) > 0 {
		group = groups[0]
	}
	return g.Redis(resolveGroup(ctx, group, gredis.DefaultGroupName))
}

// DefaultRedis 返回框架默认 Redis 实例，不读取请求 channel。
func DefaultRedis() *gredis.Redis {
	return g.Redis(gredis.DefaultGroupName)
}

func IsConfiguredChannel(_ context.Context, channel string) bool {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return false
	}
	channelsMu.RLock()
	_, ok := configuredChannels[channel]
	channelsMu.RUnlock()
	return ok
}

func resolveGroup(ctx context.Context, fallback string, defaultGroup string) string {
	fallback = strings.TrimSpace(fallback)
	channel := GetChannel(ctx)
	if channel != "" && (fallback == "" || fallback == defaultGroup) {
		return channel
	}
	if fallback != "" {
		return fallback
	}
	return defaultGroup
}

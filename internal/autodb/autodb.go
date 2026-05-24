package autodb

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type channelKey struct{}

var ctxKey = channelKey{}

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

func IsConfiguredChannel(ctx context.Context, channel string) bool {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return false
	}
	dbConfig, err := g.Cfg().Get(ctx, "database."+channel)
	if err != nil || dbConfig.IsNil() {
		return false
	}
	redisConfig, err := g.Cfg().Get(ctx, "redis."+channel)
	if err != nil || redisConfig.IsNil() {
		return false
	}
	return true
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

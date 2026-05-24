package dbcache

import (
	"context"
	"strings"
	"time"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/database/gdb"
)

// 默认缓存时长，启动时可通过 SetTTL 修改。
var defaultTTL = 5 * time.Minute

func SetTTL(d time.Duration) { defaultTTL = d }

// BuildKey 按当前 channel 和缓存键层级构建具名缓存 key。
func BuildKey(ctx context.Context, parts ...string) string {
	channel := autodb.GetChannel(ctx)
	if channel == "" {
		channel = autodb.DefaultChannelName
	}
	return strings.Join(append([]string{channel}, parts...), ":")
}

// Opt 返回一个使用默认 TTL 的 CacheOption，可选指定 name。
//
//	dao.PrfTask.Ctx(ctx).Cache(dbcache.Opt(dbcache.BuildKey(ctx, "prf_task", "id", "1001"))).Where(...)
func Opt(name ...string) gdb.CacheOption {
	o := gdb.CacheOption{Duration: defaultTTL}
	if len(name) > 0 {
		o.Name = name[0]
	}
	return o
}

// OptD 返回一个自定义 TTL 的 CacheOption，可选指定 name。
//
//	dao.PrfTask.Ctx(ctx).Cache(dbcache.OptD(30*time.Minute, dbcache.BuildKey(ctx, "prf_task", "ser", "4", "min"))).Min("id")
func OptD(duration time.Duration, name ...string) gdb.CacheOption {
	o := gdb.CacheOption{Duration: duration}
	if len(name) > 0 {
		o.Name = name[0]
	}
	return o
}

// Del 返回一个清除缓存的 CacheOption，可选指定 name。
//
//	dao.PrfTask.Ctx(ctx).Cache(dbcache.Del(dbcache.BuildKey(ctx, "prf_task", "id", "1001"))).Where(...).Update(...)
func Del(name ...string) gdb.CacheOption {
	o := gdb.CacheOption{Duration: -1}
	if len(name) > 0 {
		o.Name = name[0]
	}
	return o
}

package dbcache

import (
	"time"

	"github.com/gogf/gf/v2/database/gdb"
)

// 默认缓存时长，启动时可通过 SetTTL 修改。
var defaultTTL = 5 * time.Minute

func SetTTL(d time.Duration) { defaultTTL = d }

// Opt 返回一个使用默认 TTL 的 CacheOption，可选指定 name。
//
//	dao.User.Ctx(ctx).Cache(dbcache.Opt("user:1001")).Where(...)
func Opt(name ...string) gdb.CacheOption {
	o := gdb.CacheOption{Duration: defaultTTL}
	if len(name) > 0 {
		o.Name = name[0]
	}
	return o
}

// OptD 返回一个自定义 TTL 的 CacheOption，可选指定 name。
//
//	dao.MemConfig.Ctx(ctx).Cache(dbcache.OptD(30*time.Minute, "config:all")).All()
func OptD(duration time.Duration, name ...string) gdb.CacheOption {
	o := gdb.CacheOption{Duration: duration}
	if len(name) > 0 {
		o.Name = name[0]
	}
	return o
}

// Del 返回一个清除缓存的 CacheOption，可选指定 name。
//
//	dao.User.Ctx(ctx).Cache(dbcache.Del("user:1001")).Where(...).Update(...)
func Del(name ...string) gdb.CacheOption {
	o := gdb.CacheOption{Duration: -1}
	if len(name) > 0 {
		o.Name = name[0]
	}
	return o
}
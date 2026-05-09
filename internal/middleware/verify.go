package middleware

import (
	"math"
	"strings"

	"server_go/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// Verify 校验 login_key 有效性并防止重放攻击。
// 跳过 /user/login 接口。
func Verify(r *ghttp.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/user/login") {
		r.Middleware.Next()
		return
	}

	ctx := r.GetCtx()
	uid := r.Get("uid").Int64()
	loginKey := r.Get("login_key").String()
	platform := r.Get("platform").String()
	version := r.Get("version").String()
	tick := r.Get("tick").Int64()
	sign := r.Get("sign").String()

	if uid == 0 || loginKey == "" || platform == "" || version == "" || tick == 0 || sign == "" {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: 参数错误"})
		return
	}

	// 使用 DAO 从数据库校验 login_key
	keyVal, err := dao.UserLoginkey.Ctx(ctx).Where("uid", uid).Value("key")
	if err != nil || keyVal.IsEmpty() || keyVal.String() != loginKey {
		r.Response.WriteJsonExit(g.Map{"code": -1035, "msg": "Verify: 该账号已在其他地方登陆"})
		return
	}

	// 检查 tick 时间偏移（±1800 秒）
	now := gtime.Timestamp()
	if math.Abs(float64(now-tick)) > 1800 {
		r.Response.WriteJsonExit(g.Map{"code": -1035, "msg": "Verify: 时间校验失败"})
		return
	}

	// 通过 Redis 防重放
	redis := g.Redis()
	redisKey := "replay:" + g.NewVar(uid).String() + ":" + sign
	exists, err := redis.Do(ctx, "EXISTS", redisKey)
	if err == nil && exists.Int() > 0 {
		r.Response.WriteJsonExit(g.Map{"code": -1036, "msg": "Verify: 不能重复调用"})
		return
	}
	_, _ = redis.Do(ctx, "SET", redisKey, "1", "EX", 300)

	r.Middleware.Next()
}

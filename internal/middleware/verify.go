package middleware

import (
	"strconv"
	"strings"

	"server_go/internal/autodb"
	"server_go/internal/dao"
	"server_go/utility/dbcache"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Verify 校验 login_key，并依赖请求路径首段的 channel 切到对应库。
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

	cacheKey := dbcache.BuildKey(ctx, "login_key", "uid", strconv.FormatInt(uid, 10))
	redis := autodb.Redis(ctx)
	cached, err := redis.Do(ctx, "GET", cacheKey)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: 缓存校验失败"})
		return
	}
	if !cached.IsNil() {
		if cached.String() == loginKey {
			r.Middleware.Next()
			return
		}
		r.Response.WriteJsonExit(g.Map{"code": -1035, "msg": "Verify: 该账号已在其他地方登陆"})
		return
	}

	keyData, err := dao.UserLoginkey.Ctx(ctx).Where("uid", uid).Where("key", loginKey).One()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: 查询失败"})
		return
	}
	if keyData.IsEmpty() {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: login_key无效"})
		return
	}

	fillScript := `
local current = redis.call('GET', KEYS[1])
if current == false then
	redis.call('SETEX', KEYS[1], ARGV[2], ARGV[1])
	return 1
end
if current == ARGV[1] then
	redis.call('EXPIRE', KEYS[1], ARGV[2])
	return 1
end
return -1
`
	fillResult, err := redis.Do(ctx, "EVAL", fillScript, 1, cacheKey, loginKey, 7200)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: 缓存写入失败"})
		return
	}
	if fillResult.Int() == -1 {
		r.Response.WriteJsonExit(g.Map{"code": -1035, "msg": "Verify: 该账号已在其他地方登陆"})
		return
	}
	r.Middleware.Next()
}

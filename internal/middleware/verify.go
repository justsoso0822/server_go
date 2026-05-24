package middleware

import (
	"fmt"
	"strings"

	"server_go/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Verify校验login_key有效性并防止重放攻击。
// 跳过/user/login接口。
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

	// 1. 先查 login_key 缓存（2小时TTL），key 为 uid
	cacheKey := "login_key:uid:" + fmt.Sprintf("%d", uid)
	redis := g.Redis()
	cached, _ := redis.Do(ctx, "GET", cacheKey)
	if !cached.IsNil() {
		// 缓存命中时，验证请求中的 login_key 是否与缓存一致
		cachedKey := cached.String()
		if cachedKey == loginKey {
			r.Middleware.Next()
			return
		}
		// 缓存与请求不一致（可能被顶号），拒绝请求
		r.Response.WriteJsonExit(g.Map{"code": -1035, "msg": "Verify: 该账号已在其他地方登陆"})
		return
	}

	// 2. 查询数据库验证login_key
	keyData, err := dao.UserLoginkey.Ctx(ctx).Where("uid", uid).Where("login_key", loginKey).One()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: 查询失败"})
		return
	}
	if keyData.IsEmpty() {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: login_key无效"})
		return
	}

	// 3. 写入缓存（2小时TTL）
	redis.Do(ctx, "SETEX", cacheKey, 7200, loginKey)

	r.Middleware.Next()
}

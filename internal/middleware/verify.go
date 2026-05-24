package middleware

import (
	"fmt"
	"strings"

	"server_go/internal/autodb"
	"server_go/internal/dao"

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

	cacheKey := "login_key:uid:" + fmt.Sprintf("%d", uid)
	redis := autodb.Redis(ctx)
	cached, _ := redis.Do(ctx, "GET", cacheKey)
	if !cached.IsNil() {
		if cached.String() == loginKey {
			r.Middleware.Next()
			return
		}
		r.Response.WriteJsonExit(g.Map{"code": -1035, "msg": "Verify: 该账号已在其他地方登陆"})
		return
	}

	keyData, err := dao.UserLoginkey.Ctx(ctx).Where("uid", uid).Where("login_key", loginKey).One()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: 查询失败"})
		return
	}
	if keyData.IsEmpty() {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "Verify: login_key无效"})
		return
	}

	redis.Do(ctx, "SETEX", cacheKey, 7200, loginKey)
	r.Middleware.Next()
}

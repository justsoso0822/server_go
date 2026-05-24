package middleware

import (
	"server_go/utility/signutil"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Sign 校验请求的 HMAC-SHA256 签名。
func Sign(r *ghttp.Request) {
	errResp := g.Map{"code": -1, "msg": "非法调用"}

	// 收集所有参与签名的参数
	params := r.GetMap()
	for key := range r.GetRouterMap() {
		delete(params, key)
	}

	// 从参数或请求头获取 sign
	sign := r.Get("sign").String()
	if sign == "" {
		sign = r.GetHeader("x-sign")
	}
	if sign == "" {
		sign = r.GetHeader("x-signature")
	}
	if sign == "" {
		r.Response.WriteJsonExit(errResp)
		return
	}

	payload := signutil.BuildParams(params)

	// 从配置读取密钥
	keysVar, err := g.Cfg().Get(r.GetCtx(), "app.keys")
	if err != nil || keysVar.IsNil() {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "签名配置错误"})
		return
	}
	keys := keysVar.Strings()
	if len(keys) == 0 {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "签名配置错误"})
		return
	}

	pass := false
	for _, secret := range keys {
		if secret == "" {
			continue
		}
		computed := signutil.SHA256Hex(payload, secret)
		if computed == sign {
			pass = true
			break
		}
	}

	if !pass {
		r.Response.WriteJsonExit(errResp)
		return
	}

	r.Middleware.Next()
}

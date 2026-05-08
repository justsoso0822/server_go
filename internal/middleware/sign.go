package middleware

import (
	"server_go/utility/signutil"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Sign validates the HMAC-SHA256 signature of the request.
func Sign(r *ghttp.Request) {
	errResp := g.Map{"code": -1, "msg": "非法调用"}

	// Collect all params for signature
	params := r.GetMap()

	// Get sign from params or header
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

	// Read keys from config
	keysVar, _ := g.Cfg().Get(r.GetCtx(), "app.keys")
	if keysVar.IsNil() {
		r.Response.WriteJsonExit(errResp)
		return
	}
	keys := keysVar.Strings()

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
package middleware

import (
	"strings"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// DefaultChannel 将无渠道前缀的请求注入 default channel。
func DefaultChannel(r *ghttp.Request) {
	channel := "default"
	if !autodb.IsConfiguredChannel(r.GetCtx(), channel) {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "channel invalid"})
		return
	}

	r.SetCtx(autodb.WithChannel(r.GetCtx(), channel))
	r.Middleware.Next()
}

// Channel 将请求路径首段的 channel 注入上下文，供数据库和 Redis 路由使用。
func Channel(r *ghttp.Request) {
	channel := strings.TrimSpace(r.GetRouter("channel").String())
	if channel == "" {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "channel missing"})
		return
	}

	if !autodb.IsConfiguredChannel(r.GetCtx(), channel) {
		r.Response.WriteJsonExit(g.Map{"code": -1, "msg": "channel invalid"})
		return
	}

	r.SetCtx(autodb.WithChannel(r.GetCtx(), channel))
	r.Middleware.Next()
}

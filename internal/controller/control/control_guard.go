package control

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func ensureInternalAccess(r *ghttp.Request) bool {
	// 只允许容器内部调用，拒绝通过网关转发的请求
	forwarded := r.GetHeader("x-forwarded-for")
	if forwarded != "" {
		r.Response.Status = 404
		r.Response.WriteJson(g.Map{"ok": false})
		return false
	}
	return true
}
package middleware

import (
	"server_go/internal/logic/drainstate"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// DrainGuard 在排水阶段拒绝新请求，并维护在途请求计数。
func DrainGuard(r *ghttp.Request) {
	if drainstate.IsRejecting() {
		r.Response.Status = 503
		r.Response.WriteJsonExit(g.Map{
			"code": -1,
			"msg":  "service is draining",
		})
		return
	}

	drainstate.IncActiveRequests()
	defer drainstate.DecActiveRequests()

	r.Middleware.Next()
}

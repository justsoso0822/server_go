package middleware

import (
	"server_go/internal/runtime/drain"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// DrainGuard 在排水阶段拒绝新请求，并维护在途请求计数。
func DrainGuard(r *ghttp.Request) {
	if drain.IsRejecting() {
		r.Response.Status = 503
		r.Response.WriteJsonExit(g.Map{
			"code": -1,
			"msg":  "service is draining",
		})
		return
	}

	drain.IncActiveRequests()
	defer drain.DecActiveRequests()

	r.Middleware.Next()
}

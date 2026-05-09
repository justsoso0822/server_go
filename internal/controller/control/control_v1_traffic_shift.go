package control

import (
	"context"

	"server_go/api/control/v1"
	"server_go/internal/controller/drainstate"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) TrafficShift(ctx context.Context, req *v1.TrafficShiftReq) (res *v1.TrafficShiftRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainstate.StartTrafficShift()
	r.Response.WriteJson(g.Map{"ok": true, "state": "traffic-shift"})
	return
}

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

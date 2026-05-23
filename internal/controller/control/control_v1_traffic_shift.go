package control

import (
	"context"

	v1 "server_go/api/control/v1"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) TrafficShift(ctx context.Context, req *v1.TrafficShiftReq) (res *v1.TrafficShiftRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	state, err := service.Control().TrafficShift()
	if err != nil {
		return nil, err
	}
	r.Response.WriteJson(g.Map{"ok": true, "state": state})
	return
}

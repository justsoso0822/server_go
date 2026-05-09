package control

import (
	"context"

	"server_go/api/control/v1"
	"server_go/internal/controller/drainstate"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) RejectNew(ctx context.Context, req *v1.RejectNewReq) (res *v1.RejectNewRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainstate.StartRejectNew()
	r.Response.WriteJson(g.Map{"ok": true, "state": "reject-new-requests"})
	return
}

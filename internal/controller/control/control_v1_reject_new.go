package control

import (
	"context"

	v1 "server_go/api/control/v1"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) RejectNew(ctx context.Context, req *v1.RejectNewReq) (res *v1.RejectNewRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	state, err := service.Control().RejectNew()
	if err != nil {
		return nil, err
	}
	r.Response.WriteJson(g.Map{"ok": true, "state": state})
	return
}

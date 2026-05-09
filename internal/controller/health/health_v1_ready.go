package health

import (
	"context"

	"server_go/api/health/v1"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) Ready(ctx context.Context, req *v1.ReadyReq) (res *v1.ReadyRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{"ok": true})
	return
}

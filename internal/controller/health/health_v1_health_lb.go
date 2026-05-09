package health

import (
	"context"

	"server_go/api/health/v1"
	"server_go/internal/controller/drainstate"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) HealthLb(ctx context.Context, req *v1.HealthLbReq) (res *v1.HealthLbRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if drainstate.IsTrafficShift() {
		r.Response.Status = 503
		r.Response.WriteJson(g.Map{"status": "draining"})
	} else {
		r.Response.WriteJson(g.Map{"status": "ok"})
	}
	return
}

package health

import (
	"context"
	"os"
	"time"

	"server_go/api/health/v1"
	"server_go/internal/controller/drainstate"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) HealthDetail(ctx context.Context, req *v1.HealthDetailReq) (res *v1.HealthDetailRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":    "ok",
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format("2006/01/02 15:04:05"),
		"draining":  drainstate.IsTrafficShift(),
		"rejecting": drainstate.IsRejecting(),
	})
	return
}

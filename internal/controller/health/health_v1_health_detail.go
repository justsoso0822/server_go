package health

import (
	"context"
	"os"
	"time"

	"server_go/api/health/v1"
	"server_go/internal/controller/drainstate"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerV1) HealthDetail(ctx context.Context, req *v1.HealthDetailReq) (res *v1.HealthDetailRes, err error) {
	color := genv.Get("APP_COLOR", "unknown").String()

	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":         "ok",
		"color":          color,
		"pid":            os.Getpid(),
		"uptime":         int(time.Since(startTime).Seconds()),
		"timestamp":      gtime.Now().Format("Y/m/d H:i:s"),
		"draining":       drainstate.IsTrafficShift(),
		"rejecting":      drainstate.IsRejecting(),
		"activeRequests": drainstate.GetActiveRequests(),
	})
	return
}

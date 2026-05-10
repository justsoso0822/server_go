package health

import (
	"context"
	"os"
	"time"

	"server_go/api/health/v1"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerV1) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	color := genv.Get("APP_COLOR", "unknown").String()

	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":    "ok",
		"color":     color,
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": gtime.Now().Format("Y/m/d H:i:s"),
	})
	return
}

package gamelog

import (
	"context"
	"math"

	"server_go/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// TraceRes records a resource change asynchronously.
func TraceRes(ctx context.Context, uid int64, old, now int64, resName, reason string) {
	if uid == 0 {
		return
	}
	num := now - old
	if num == 0 {
		return
	}
	if num > 0 {
		resName = "+" + resName
	} else {
		resName = "-" + resName
	}
	absNum := int64(math.Abs(float64(num)))
	bgCtx := gctx.NeverDone(ctx)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "TraceRes panic: %v", r)
			}
		}()
		_, _ = dao.LogTrace.Ctx(bgCtx).Data(g.Map{
			"uid": uid, "type": resName, "num": absNum,
			"before": old, "after": now, "reason": reason,
		}).Insert()
	}()
}

// Log records a message asynchronously.
func Log(ctx context.Context, uid int64, msg string) {
	bgCtx := gctx.NeverDone(ctx)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "Log panic: %v", r)
			}
		}()
		_, _ = dao.LogMsg.Ctx(bgCtx).Data(g.Map{"uid": uid, "msg": msg}).Insert()
	}()
}
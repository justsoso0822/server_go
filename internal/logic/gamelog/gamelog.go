package gamelog

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/frame/g"
)

// TraceRes records a resource change asynchronously.
func TraceRes(ctx context.Context, uid, old, now int64, resName, reason string) {
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

	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "TraceRes panic: %v", r)
			}
		}()
		_, _ = g.DB().Exec(context.Background(),
			"INSERT INTO `_log_trace` (`uid`, `type`, `num`, `before`, `after`, `reason`) VALUES (?, ?, ?, ?, ?, ?)",
			uid, resName, absNum, old, now, reason,
		)
	}()
}

// Log records a message asynchronously.
func Log(ctx context.Context, uid int64, msg string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(context.Background(), "Log panic: %v", r)
			}
		}()
		_, _ = g.DB().Exec(context.Background(),
			"INSERT INTO `_log_msg` (`uid`, `msg`) VALUES (?, ?)",
			uid, msg,
		)
	}()
}
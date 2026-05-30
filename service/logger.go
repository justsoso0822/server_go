package service

import (
	"context"
	"math"

	"server_gin/dao"
	"server_gin/dao/model"
	"server_gin/tools/autodb"
)

// TraceRes 异步记录资源变化到 log_trace。
func TraceRes(ctx context.Context, uid int64, old, now int64, resName, reason string) {
	if uid == 0 {
		return
	}
	num := now - old
	if num == 0 {
		return
	}
	label := resName
	if num > 0 {
		label = "+" + resName
	} else {
		label = "-" + resName
	}
	absNum := int64(math.Abs(float64(num)))
	bgCtx := autodb.BackgroundWithChannel(ctx)

	go func() {
		defer func() { recover() }()
		_ = dao.InsertLogTrace(bgCtx, &model.LogTrace{
			UID:    int32(uid),
			Type:   label,
			Num:    int32(absNum),
			Before: int32(old),
			After:  int32(now),
			Reason: reason,
		})
	}()
}

// LogMsg 异步记录消息到 log_msg。
func LogMsg(ctx context.Context, uid int64, msg string) {
	bgCtx := autodb.BackgroundWithChannel(ctx)
	go func() {
		defer func() { recover() }()
		_ = dao.InsertLogMsg(bgCtx, &model.LogMsg{UID: int32(uid), Msg: msg})
	}()
}

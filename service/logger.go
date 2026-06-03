package service

import (
	"context"
	"time"

	"server_go/dao"
	"server_go/dao/model"
	"server_go/tools/autodb"

	"go.uber.org/zap"
)

const asyncLogTimeout = 10 * time.Second

func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

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
	absNum := absInt64(num)
	bgCtx := autodb.BackgroundWithChannel(ctx)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logAsyncPanic(bgCtx, "async trace resource panic", r,
					zap.Int64("uid", uid),
					zap.String("res_name", resName),
					zap.String("reason", reason),
					zap.Int64("delta", num),
				)
			}
		}()
		tCtx, cancel := context.WithTimeout(bgCtx, asyncLogTimeout)
		defer cancel()
		if err := dao.InsertLogTrace(tCtx, &model.LogTrace{
			UID:       uid,
			Type:      label,
			Num:       int32(absNum),
			Before:    int32(old),
			After:     int32(now),
			Reason:    reason,
			RequestID: autodb.GetRequestID(bgCtx),
		}); err != nil {
			logAsyncError(bgCtx, "async trace resource failed", err,
				zap.Int64("uid", uid),
				zap.String("res_name", resName),
				zap.String("reason", reason),
				zap.Int64("delta", num),
			)
		}
	}()
}

func LogMsg(ctx context.Context, uid int64, msg string) {
	bgCtx := autodb.BackgroundWithChannel(ctx)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logAsyncPanic(bgCtx, "async message log panic", r,
					zap.Int64("uid", uid),
					zap.String("log_msg", msg),
				)
			}
		}()
		tCtx, cancel := context.WithTimeout(bgCtx, asyncLogTimeout)
		defer cancel()
		if err := dao.InsertLogMsg(tCtx, &model.LogMsg{
			UID:       uid,
			Msg:       msg,
			RequestID: autodb.GetRequestID(bgCtx),
		}); err != nil {
			logAsyncError(bgCtx, "async message log failed", err,
				zap.Int64("uid", uid),
				zap.String("log_msg", msg),
			)
		}
	}()
}

func logAsyncError(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = appendLogContextFields(ctx, fields...)
	fields = append(fields, zap.Error(err))
	zap.L().Error(msg, fields...)
}

func logAsyncPanic(ctx context.Context, msg string, recovered any, fields ...zap.Field) {
	fields = appendLogContextFields(ctx, fields...)
	fields = append(fields, zap.Any("panic", recovered))
	zap.L().Error(msg, fields...)
}

func appendLogContextFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	if requestID := autodb.GetRequestID(ctx); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if channel := autodb.GetChannel(ctx); channel != "" {
		fields = append(fields, zap.String("channel", channel))
	}
	if uid := autodb.GetLogUID(ctx); uid != 0 {
		fields = append(fields, zap.Int64("uid", uid))
	}
	if openid := autodb.GetLogOpenID(ctx); openid != "" {
		fields = append(fields, zap.String("openid", openid))
	}
	return fields
}

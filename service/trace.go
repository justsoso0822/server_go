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

// 资源流水只关心变化量大小，正负号已经体现在 Type 的 +/- 前缀里。
func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// 记录资源变动流水，调用方不等待写库完成。
// 这类日志是审计/排查辅助数据，不应该让日志库抖动拖慢发放资源的主流程。
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
		// DAO 层会通过 autodb.DB(tCtx) 取带 channel 的 GORM 连接，
		// 所以异步日志仍会写到当前请求所属渠道的数据库。
		if err := dao.InsertLogTrace(tCtx, &model.LogTrace{
			UID:       uid,
			Type:      label,
			Num:       absNum,
			Before:    old,
			After:     now,
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

// 写入一条普通业务消息日志，同样采用异步方式。
// 当前主要用于记录主流程中的非致命异常，例如资源更新失败的业务上下文。
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

// 异步写库失败时转写 zap 错误日志，避免错误静默丢失。
func logAsyncError(ctx context.Context, msg string, err error, fields ...zap.Field) {
	fields = appendLogContextFields(ctx, fields...)
	fields = append(fields, zap.Error(err))
	// zap.L() 使用 bootstrap 阶段 ReplaceGlobals 设置的全局 logger；
	// 这样不从参数层层传 logger，异步工具函数也能保持统一日志出口。
	zap.L().Error(msg, fields...)
}

// 后台 goroutine panic 时统一收敛到 zap，避免进程因为辅助日志任务崩溃。
func logAsyncPanic(ctx context.Context, msg string, recovered any, fields ...zap.Field) {
	fields = appendLogContextFields(ctx, fields...)
	fields = append(fields, zap.Any("panic", recovered))
	zap.L().Error(msg, fields...)
}

// 把请求链路字段补到异步日志里，让 zap 日志、GORM SQL 日志和业务表日志能互相关联。
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

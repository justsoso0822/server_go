package autodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const defaultGormSlowThreshold = 300 * time.Millisecond

type slogGormLogger struct {
	log           *slog.Logger
	level         gormlogger.LogLevel
	channel       string
	slowThreshold time.Duration
}

func newGormLogger(log *slog.Logger, channel string, debug bool, slowThresholdMS int) gormlogger.Interface {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}

	level := gormlogger.Warn
	if debug {
		// Info 级别会输出所有 SQL，适合本地定位 ORM 生成的语句；线上通常只开 Warn，
		// 只记录慢查询和错误，避免日志量暴涨。
		level = gormlogger.Info
	}
	slowThreshold := defaultGormSlowThreshold
	if slowThresholdMS > 0 {
		slowThreshold = time.Duration(slowThresholdMS) * time.Millisecond
	}

	return &slogGormLogger{
		log:           log.With("component", "gorm"),
		level:         level,
		channel:       channel,
		slowThreshold: slowThreshold,
	}
}

func (l *slogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	// GORM 可能在会话级临时调整日志级别。返回副本可以避免影响共享 logger 的默认级别。
	next := *l
	next.level = level
	return &next
}

func (l *slogGormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level < gormlogger.Info {
		return
	}
	l.log.InfoContext(ctx, fmt.Sprintf(msg, args...), slog.String("db_channel", l.channel))
}

func (l *slogGormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level < gormlogger.Warn {
		return
	}
	l.log.WarnContext(ctx, fmt.Sprintf(msg, args...), slog.String("db_channel", l.channel))
}

func (l *slogGormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level < gormlogger.Error {
		return
	}
	l.log.ErrorContext(ctx, fmt.Sprintf(msg, args...), slog.String("db_channel", l.channel))
}

func (l *slogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		// ErrRecordNotFound 在业务查询里很常见，作为错误日志会制造大量噪音。
		sql, rows := fc()
		l.log.LogAttrs(ctx, slog.LevelError, "gorm query error", l.attrs(ctx, sql, rows, elapsed, slog.Any("error", err))...)
	case elapsed > l.slowThreshold && l.level >= gormlogger.Warn:
		sql, rows := fc()
		l.log.LogAttrs(ctx, slog.LevelWarn, "gorm slow query", l.attrs(ctx, sql, rows, elapsed)...)
	case l.level >= gormlogger.Info:
		// fc 会格式化 SQL 和 rows，只有确定要写日志时才调用，减少热路径开销。
		sql, rows := fc()
		l.log.LogAttrs(ctx, slog.LevelInfo, "gorm query", l.attrs(ctx, sql, rows, elapsed)...)
	}
}

func (l *slogGormLogger) attrs(ctx context.Context, sql string, rows int64, elapsed time.Duration, extra ...slog.Attr) []slog.Attr {
	// 这些字段来自 autodb.DB(ctx).WithContext(ctx)，能把一条 SQL 串回 HTTP 请求和用户。
	attrs := []slog.Attr{
		slog.String("request_id", GetRequestID(ctx)),
		slog.String("db_channel", l.channel),
		slog.String("sql", sql),
		slog.Int64("rows", rows),
		slog.Float64("elapsed_ms", float64(elapsed.Microseconds())/1000),
	}
	if uid := GetLogUID(ctx); uid != 0 {
		attrs = append(attrs, slog.Int64("uid", uid))
	}
	if openid := GetLogOpenID(ctx); openid != "" {
		attrs = append(attrs, slog.String("openid", openid))
	}
	return append(attrs, extra...)
}

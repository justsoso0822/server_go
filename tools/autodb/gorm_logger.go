package autodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const defaultGormSlowThreshold = 300 * time.Millisecond

type zapGormLogger struct {
	log           *zap.Logger
	level         gormlogger.LogLevel
	channel       string
	slowThreshold time.Duration
}

func newGormLogger(log *zap.Logger, channel string, debug bool, slowThresholdMS int) gormlogger.Interface {
	if log == nil {
		log = zap.NewNop()
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

	return &zapGormLogger{
		log:           log.Named("gorm"),
		level:         level,
		channel:       channel,
		slowThreshold: slowThreshold,
	}
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	// GORM 可能在会话级临时调整日志级别。返回副本可以避免影响共享 logger 的默认级别。
	next := *l
	next.level = level
	return &next
}

func (l *zapGormLogger) Info(_ context.Context, msg string, args ...any) {
	if l.level < gormlogger.Info {
		return
	}
	l.log.Info(fmt.Sprintf(msg, args...), zap.String("db_channel", l.channel))
}

func (l *zapGormLogger) Warn(_ context.Context, msg string, args ...any) {
	if l.level < gormlogger.Warn {
		return
	}
	l.log.Warn(fmt.Sprintf(msg, args...), zap.String("db_channel", l.channel))
}

func (l *zapGormLogger) Error(_ context.Context, msg string, args ...any) {
	if l.level < gormlogger.Error {
		return
	}
	l.log.Error(fmt.Sprintf(msg, args...), zap.String("db_channel", l.channel))
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		// ErrRecordNotFound 在业务查询里很常见，作为错误日志会制造大量噪音。
		sql, rows := fc()
		l.log.Error("gorm query error", l.fields(ctx, sql, rows, elapsed, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.level >= gormlogger.Warn:
		sql, rows := fc()
		l.log.Warn("gorm slow query", l.fields(ctx, sql, rows, elapsed)...)
	case l.level >= gormlogger.Info:
		// fc 会格式化 SQL 和 rows，只有确定要写日志时才调用，减少热路径开销。
		sql, rows := fc()
		l.log.Info("gorm query", l.fields(ctx, sql, rows, elapsed)...)
	}
}

func (l *zapGormLogger) fields(ctx context.Context, sql string, rows int64, elapsed time.Duration, extra ...zap.Field) []zap.Field {
	// 这些字段来自 autodb.DB(ctx).WithContext(ctx)，能把一条 SQL 串回 HTTP 请求和用户。
	fields := []zap.Field{
		zap.String("request_id", GetRequestID(ctx)),
		zap.String("db_channel", l.channel),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Float64("elapsed_ms", float64(elapsed.Microseconds())/1000),
	}
	if uid := GetLogUID(ctx); uid != 0 {
		fields = append(fields, zap.Int64("uid", uid))
	}
	if openid := GetLogOpenID(ctx); openid != "" {
		fields = append(fields, zap.String("openid", openid))
	}
	return append(fields, extra...)
}

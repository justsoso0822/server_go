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

func (l *zapGormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		l.log.Error("gorm query error", l.fields(sql, rows, elapsed, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.level >= gormlogger.Warn:
		sql, rows := fc()
		l.log.Warn("gorm slow query", l.fields(sql, rows, elapsed)...)
	case l.level >= gormlogger.Info:
		sql, rows := fc()
		l.log.Info("gorm query", l.fields(sql, rows, elapsed)...)
	}
}

func (l *zapGormLogger) fields(sql string, rows int64, elapsed time.Duration, extra ...zap.Field) []zap.Field {
	fields := []zap.Field{
		zap.String("db_channel", l.channel),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Float64("elapsed_ms", float64(elapsed.Microseconds())/1000),
	}
	return append(fields, extra...)
}

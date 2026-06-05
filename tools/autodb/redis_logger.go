package autodb

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// zapRedisLogger 适配 go-redis 的 internal.Logging 接口,
// 将 go-redis 内部的底层告警(连接池异常、重连失败等)转发到 zap。
// go-redis 默认通过标准库 log 输出到 stderr,会以裸文本混入 JSON 日志流,
// 接管后统一格式且不丢底层告警。
type zapRedisLogger struct {
	log *zap.Logger
}

func newRedisLogger(log *zap.Logger) *zapRedisLogger {
	if log == nil {
		log = zap.NewNop()
	}
	return &zapRedisLogger{log: log.Named("redis")}
}

// Printf 是 go-redis internal.Logging 接口的唯一方法。
// go-redis 内部仅用它输出告警级别的信息,统一映射为 zap Warn。
func (l *zapRedisLogger) Printf(_ context.Context, format string, v ...any) {
	l.log.Warn(fmt.Sprintf(format, v...))
}

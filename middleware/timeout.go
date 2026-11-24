package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// DefaultRequestTimeout 是请求链路的兜底上限。
const DefaultRequestTimeout = 2 * time.Second

// RequestTimeout 给每个请求的 context 挂一个截止时间。
//
// 它不是强制中断：Go 的 context 只在协作点生效，handler 里如果完全不检查 ctx，
// 这段代码也拦不住它。但下游的 Redis / MySQL 客户端都认这个 deadline——
// go-redis 和 gorm 都会在到点时提前返回 context deadline exceeded，
// 请求不会因为一个慢依赖挂到 HTTP 层自己超时。
func RequestTimeout(d time.Duration) gin.HandlerFunc {
	if d <= 0 {
		d = DefaultRequestTimeout
	}
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		// cancel 必须在整个请求结束后调用，否则子 context 会挂到父 context 结束才释放。
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

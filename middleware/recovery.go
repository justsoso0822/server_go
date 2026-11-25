package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"

	"app/tools/autodb"

	"github.com/gin-gonic/gin"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}

	// Gin 默认 Recovery 会把 panic 写到 DefaultErrorWriter。这里把 writer 指向 io.Discard，
	// 再手动写结构化日志，避免线上 JSON 日志中混入多行文本堆栈。
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		ctx := c.Request.Context()
		attrs := []slog.Attr{
			slog.String("request_id", autodb.GetRequestID(ctx)),
			slog.String("error", fmt.Sprint(recovered)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Any("query_keys", queryKeys(c.Request.URL.Query())),
			slog.String("client_ip", c.ClientIP()),
			slog.String("stack", string(debug.Stack())),
		}
		if uid := autodb.GetLogUID(ctx); uid != 0 {
			attrs = append(attrs, slog.Int64("uid", uid))
		}
		if openid := autodb.GetLogOpenID(ctx); openid != "" {
			attrs = append(attrs, slog.String("openid", openid))
		}
		log.LogAttrs(ctx, slog.LevelError, "panic recovered", attrs...)

		if c.Writer.Written() {
			c.Abort()
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "internal server error",
		})
	})
}

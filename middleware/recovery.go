package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"server_go/tools/autodb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	if log == nil {
		log = zap.NewNop()
	}

	// Gin 默认 Recovery 会把 panic 写到 DefaultErrorWriter。这里把 writer 指向 io.Discard，
	// 再手动写 zap，避免线上 JSON 日志中混入多行文本堆栈。
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		ctx := c.Request.Context()
		fields := []zap.Field{
			zap.String("request_id", autodb.GetRequestID(c.Request.Context())),
			zap.String("error", fmt.Sprint(recovered)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Strings("query_keys", queryKeys(c.Request.URL.Query())),
			zap.String("client_ip", c.ClientIP()),
			zap.String("stack", string(debug.Stack())),
		}
		if uid := autodb.GetLogUID(ctx); uid != 0 {
			fields = append(fields, zap.Int64("uid", uid))
		}
		if openid := autodb.GetLogOpenID(ctx); openid != "" {
			fields = append(fields, zap.String("openid", openid))
		}
		log.Error("panic recovered", fields...)

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

package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"server_gin/tools/autodb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery logs panics through zap so non-local container logs stay single-line JSON.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	if log == nil {
		log = zap.NewNop()
	}

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

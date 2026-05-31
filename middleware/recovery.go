package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery logs panics through zap so non-local container logs stay single-line JSON.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	if log == nil {
		log = zap.NewNop()
	}

	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		log.Error("panic recovered",
			zap.String("error", fmt.Sprint(recovered)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("client_ip", c.ClientIP()),
			zap.String("stack", string(debug.Stack())),
		)

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

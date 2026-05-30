package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// TestEnvGuard 限制 /test 组只在 local/test 环境可访问。
func TestEnvGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
		if env == "" {
			env = "local"
		}
		if env != "local" && env != "test" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": -1,
				"msg":  "test endpoints are disabled in current environment",
				"env":  env,
			})
			return
		}
		c.Next()
	}
}

// InternalOnly 只允许容器内部直连调用，拒绝经网关转发的请求。
func InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("x-forwarded-for") != "" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"ok": false})
			return
		}
		c.Next()
	}
}

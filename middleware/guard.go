package middleware

import (
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var forwardedHeaders = [...]string{
	"Forwarded",
	"X-Forwarded-For",
	"X-Forwarded-Host",
	"X-Forwarded-Proto",
	"X-Real-IP",
}

// 限制 /test 组只在 local/test 环境可访问。
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

// 只允许容器内部直连调用，拒绝经网关转发的请求。
func InternalOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if hasForwardedHeader(c) || !isInternalRemoteAddr(c.Request.RemoteAddr) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"ok": false})
			return
		}
		c.Next()
	}
}

func hasForwardedHeader(c *gin.Context) bool {
	for _, header := range forwardedHeaders {
		if c.GetHeader(header) != "" {
			return true
		}
	}
	return false
}

func isInternalRemoteAddr(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(strings.TrimSpace(remoteAddr))
	if err != nil {
		host = strings.TrimSpace(remoteAddr)
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	return addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast()
}

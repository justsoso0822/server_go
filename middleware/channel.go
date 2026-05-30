package middleware

import (
	"net/http"
	"strings"

	"server_gin/tools/autodb"

	"github.com/gin-gonic/gin"
)

// DefaultChannel 将无渠道前缀的请求注入 default channel。
func DefaultChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !autodb.IsConfiguredChannel(autodb.DefaultChannelName) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "channel invalid"})
			return
		}
		ctx := autodb.WithChannel(c.Request.Context(), autodb.DefaultChannelName)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Channel 从路由参数取 channel 名并注入 context。
func Channel() gin.HandlerFunc {
	return func(c *gin.Context) {
		channel := strings.TrimSpace(c.Param("channel"))
		if channel == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "channel missing"})
			return
		}
		if !autodb.IsConfiguredChannel(channel) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "channel invalid"})
			return
		}
		ctx := autodb.WithChannel(c.Request.Context(), channel)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

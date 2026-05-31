package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"server_go/state"
	"server_go/tools/autodb"

	"github.com/gin-gonic/gin"
)

const replayWindow = 5 * time.Minute

// ReplayGuard 校验 tick 有效期，并拒绝已使用过的签名请求。
func ReplayGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := collectRequestParams(c)
		sign := requestSign(c, params)
		tick := paramString(params, "tick")
		if sign == "" || tick == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "非法调用"})
			return
		}

		ts, err := parseTick(tick)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "请求已过期"})
			return
		}
		now := time.Now()
		if ts.Before(now.Add(-replayWindow)) || ts.After(now.Add(replayWindow)) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "请求已过期"})
			return
		}

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		key := state.BuildKey(c.Request.Context(), "replay", c.Request.Method, path, sign)
		ttl := time.Until(ts.Add(replayWindow))
		if ttl <= 0 {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "请求已过期"})
			return
		}
		ok, err := autodb.Redis(c.Request.Context()).SetNX(c.Request.Context(), key, "1", ttl).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "防重放校验失败"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "请求重复"})
			return
		}

		c.Next()
	}
}

func parseTick(tick string) (time.Time, error) {
	raw := strings.TrimSpace(tick)
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	if n > 1_000_000_000_000 {
		return time.UnixMilli(n), nil
	}
	return time.Unix(n, 0), nil
}

func paramString(params map[string]interface{}, key string) string {
	if v, ok := params[key].(string); ok {
		return v
	}
	return ""
}

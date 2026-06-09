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

func ReplayGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params map[string]interface{}
		var sign string
		if v, ok := c.Get("_params"); ok {
			params = v.(map[string]interface{})
		} else {
			params = collectRequestParams(c)
		}
		if v, ok := c.Get("_sign"); ok {
			sign = v.(string)
		} else {
			sign = requestSign(c, params)
		}
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
		// FullPath 是 Gin 匹配后的模板路径，如 /api/bag/get_bag/:chapter。
		key := state.BuildKey(c.Request.Context(), "replay", c.Request.Method, path, sign)
		ttl := time.Until(ts.Add(replayWindow))
		if ttl <= 0 {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "请求已过期"})
			return
		}
		// Redis SET NX EX/PX 是原子操作：只有第一次请求能写入成功，重复请求会返回 false。
		// 扩展知识：如果要防跨区域重放，需要保证所有实例使用同一个 Redis 或同等一致性的存储。
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

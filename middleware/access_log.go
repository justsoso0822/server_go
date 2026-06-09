package middleware

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"server_go/tools/autodb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 记录 HTTP 请求日志，自动过滤健康检查和控制接口。
func AccessLog(log *zap.Logger) gin.HandlerFunc {
	if log == nil {
		log = zap.NewNop()
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		if shouldSkipAccessLog(path) {
			return
		}

		ctx := c.Request.Context()
		status := c.Writer.Status()
		fields := []zap.Field{
			zap.String("request_id", autodb.GetRequestID(ctx)),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("route", c.FullPath()),
			zap.Strings("query_keys", queryKeys(c.Request.URL.Query())),
			zap.Int("status", status),
			zap.Int("body_size", c.Writer.Size()),
			zap.Float64("latency_ms", float64(time.Since(start).Microseconds())/1000),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("channel", autodb.GetChannel(ctx)),
		}
		if uid := autodb.GetLogUID(ctx); uid != 0 {
			fields = append(fields, zap.Int64("uid", uid))
		} else if uid := requestParam(c, "uid"); uid != "" {
			fields = appendUIDField(fields, uid)
		}
		if openid := autodb.GetLogOpenID(ctx); openid != "" {
			fields = append(fields, zap.String("openid", openid))
		}
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case status >= 500:
			log.Error("http access", fields...)
		case status >= 400:
			log.Warn("http access", fields...)
		default:
			log.Info("http access", fields...)
		}
	}
}

func appendUIDField(fields []zap.Field, uid string) []zap.Field {
	if n, err := strconv.ParseInt(uid, 10, 64); err == nil {
		return append(fields, zap.Int64("uid", n))
	}
	return append(fields, zap.String("uid", uid))
}

func shouldSkipAccessLog(path string) bool {
	path = strings.TrimRight(path, "/")
	return path == "/health" || strings.HasPrefix(path, "/health/") ||
		strings.HasPrefix(path, "/internal/control/")
}

func requestParam(c *gin.Context, key string) string {
	if v := c.Query(key); v != "" {
		return v
	}
	if c.Request.PostForm == nil {
		return ""
	}
	if vals := c.Request.PostForm[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func queryKeys(values url.Values) []string {
	if len(values) == 0 {
		return nil
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

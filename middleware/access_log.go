package middleware

import (
	"log/slog"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"app/tools/autodb"

	"github.com/gin-gonic/gin"
)

// 记录 HTTP 请求日志，自动过滤健康检查和控制接口。
func AccessLog(log *slog.Logger) gin.HandlerFunc {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
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
		attrs := []slog.Attr{
			slog.String("request_id", autodb.GetRequestID(ctx)),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("route", c.FullPath()),
			slog.Any("query_keys", queryKeys(c.Request.URL.Query())),
			slog.Int("status", status),
			slog.Int("body_size", c.Writer.Size()),
			slog.Float64("latency_ms", float64(time.Since(start).Microseconds())/1000),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("channel", autodb.GetChannel(ctx)),
		}
		if uid := autodb.GetLogUID(ctx); uid != 0 {
			attrs = append(attrs, slog.Int64("uid", uid))
		} else if uid := requestParam(c, "uid"); uid != "" {
			attrs = appendUIDAttr(attrs, uid)
		}
		if openid := autodb.GetLogOpenID(ctx); openid != "" {
			attrs = append(attrs, slog.String("openid", openid))
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		switch {
		case status >= 500:
			log.LogAttrs(ctx, slog.LevelError, "http access", attrs...)
		case status >= 400:
			log.LogAttrs(ctx, slog.LevelWarn, "http access", attrs...)
		default:
			log.LogAttrs(ctx, slog.LevelInfo, "http access", attrs...)
		}
	}
}

func appendUIDAttr(attrs []slog.Attr, uid string) []slog.Attr {
	if n, err := strconv.ParseInt(uid, 10, 64); err == nil {
		return append(attrs, slog.Int64("uid", n))
	}
	return append(attrs, slog.String("uid", uid))
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

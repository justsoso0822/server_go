package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 预检请求缺少 Access-Control-Request-Headers 时使用的兜底允许头。
const (
	corsAllowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowHeaders = "Content-Type, Authorization, X-Requested-With, X-Request-Id"
	corsMaxAge       = "86400"
)

// CORS 处理浏览器跨域请求。
//
// H5 调试时页面由 LayaAir 内置服务器提供（如 http://localhost:18090），
// API 走 Traefik 网关（如 http://localhost:17001），两者源不同；
// 服务端不返回 Access-Control-Allow-Origin 时，浏览器会直接阻断响应，
// 前端只能看到 "Request failed Status: 0"。
//
// 放行规则（allowedOrigins 来自配置 app.allowedOrigins）：
//   - 含 "*"：回显任意 Origin；
//   - 非空：只回显名单内的 Origin，其余不返回任何 CORS 头；
//   - 为空：local/test 环境回显任意 Origin，方便局域网和 IDE 内置服务器联调；
//     生产环境视为不放行跨域。
//
// 必须挂在验签中间件（middleware.Sign）之前注册：预检 OPTIONS 不带 sign 参数，
// 若先进入验签会被判定为“非法调用”并返回 200，浏览器侧依旧认为是 CORS 失败。
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowAll := false
	whitelist := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		origin = strings.ToLower(strings.TrimSpace(origin))
		switch {
		case origin == "":
			continue
		case origin == "*":
			allowAll = true
		default:
			whitelist[origin] = struct{}{}
		}
	}
	// 白名单为空时按环境兜底：本地和测试环境宽松放行，生产环境必须显式配置。
	if !allowAll && len(whitelist) == 0 && isLocalLikeEnv() {
		allowAll = true
	}

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		if origin == "" {
			// 非浏览器跨域请求（小程序、原生、curl 等）不需要 CORS 头。
			c.Next()
			return
		}

		if !allowAll {
			if _, ok := whitelist[strings.ToLower(origin)]; !ok {
				if c.Request.Method == http.MethodOptions {
					// 不在白名单的预检直接拒绝，不再进入后续中间件。
					c.AbortWithStatus(http.StatusNoContent)
					return
				}
				c.Next()
				return
			}
		}

		header := c.Writer.Header()
		// 回显具体 Origin（而不是 *），因为业务请求携带 query 签名，
		// 后续若启用 Cookie 也无需修改这里。
		header.Set("Access-Control-Allow-Origin", origin)
		// 回显 Origin 的响应与请求源相关，必须声明 Vary，避免网关/CDN 把某个源的响应头复用到其它源。
		header.Add("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			header.Set("Access-Control-Allow-Methods", corsAllowMethods)
			// 浏览器会声明它想使用的请求头，直接回显即可覆盖 Content-Type 与自定义头。
			if requestHeaders := strings.TrimSpace(c.GetHeader("Access-Control-Request-Headers")); requestHeaders != "" {
				header.Set("Access-Control-Allow-Headers", requestHeaders)
			} else {
				header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
			}
			header.Set("Access-Control-Max-Age", corsMaxAge)
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isLocalLikeEnv 判断当前是否本地/测试环境；未设置 APP_ENV 时按本地处理，与 configureGinMode 保持一致。
func isLocalLikeEnv() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "", "local", "test":
		return true
	default:
		return false
	}
}

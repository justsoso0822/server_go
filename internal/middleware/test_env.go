package middleware

import (
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const defaultTestEnv = "local"

var allowedTestEnvs = map[string]struct{}{
	"local": {},
	"test":  {},
}

// TestEnvGuard 限制 /test 分组只能在 local/test 环境访问。
func TestEnvGuard(r *ghttp.Request) {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env == "" {
		env = defaultTestEnv
	}
	if _, ok := allowedTestEnvs[env]; !ok {
		r.Response.Status = 403
		r.Response.WriteJsonExit(g.Map{
			"code": -1,
			"msg":  "test endpoints are disabled in current environment",
			"env":  env,
		})
		return
	}

	r.Middleware.Next()
}
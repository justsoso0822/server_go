package middleware

import (
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Response wraps handler output with {code: 0, data: ...} if needed.
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.Status >= 400 {
		return
	}

	buf := r.Response.Buffer()
	if len(buf) == 0 {
		return
	}

	var result interface{}
	if err := json.Unmarshal(buf, &result); err != nil {
		r.Response.ClearBuffer()
		r.Response.WriteJson(g.Map{"code": 0, "data": string(buf)})
		return
	}

	if m, ok := result.(map[string]interface{}); ok {
		if _, hasCode := m["code"]; !hasCode {
			m["code"] = 0
		}
		r.Response.ClearBuffer()
		r.Response.WriteJson(m)
		return
	}

	r.Response.ClearBuffer()
	r.Response.WriteJson(g.Map{"code": 0, "data": result})
}
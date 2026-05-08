package middleware

import (
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Response writes controller return values using GoFrame's handler response lifecycle.
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.Status >= 400 || len(r.Response.Buffer()) > 0 {
		return
	}

	if err := r.GetError(); err != nil {
		r.Response.WriteJson(g.Map{"code": -1, "msg": err.Error()})
		return
	}

	result := r.GetHandlerResponse()
	if result == nil {
		return
	}

	if m, ok := result.(g.Map); ok {
		if _, hasCode := m["code"]; !hasCode {
			m["code"] = 0
		}
		r.Response.WriteJson(m)
		return
	}

	var object map[string]interface{}
	b, err := json.Marshal(result)
	if err == nil && json.Unmarshal(b, &object) == nil && object != nil {
		if _, hasCode := object["code"]; !hasCode {
			object["code"] = 0
		}
		r.Response.WriteJson(object)
		return
	}

	r.Response.WriteJson(g.Map{"code": 0, "data": result})
}

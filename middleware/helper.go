package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func firstStr(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func collectRequestParams(c *gin.Context) map[string]interface{} {
	params := map[string]interface{}{}
	for k, v := range c.Request.URL.Query() {
		params[k] = strings.Join(v, ",")
	}
	if c.Request.Method == "POST" {
		_ = c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			params[k] = strings.Join(v, ",")
		}
	}
	return params
}

func requestSign(c *gin.Context, params map[string]interface{}) string {
	if s, ok := params["sign"].(string); ok {
		return s
	}
	return ""
}

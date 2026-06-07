package middleware

import (
	"bytes"
	"encoding/json"
	"io"
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
		if strings.Contains(c.ContentType(), "application/json") {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// 还原 body，供后续 handler 的 ShouldBindJSON 使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
				var jsonParams map[string]interface{}
				if json.Unmarshal(body, &jsonParams) == nil {
					for k, v := range jsonParams {
						params[k] = v
					}
				}
			}
		} else {
			_ = c.Request.ParseForm()
			for k, v := range c.Request.PostForm {
				params[k] = strings.Join(v, ",")
			}
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

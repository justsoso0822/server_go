package middleware

import (
	"net/http"
	"strings"

	"server_gin/config"
	"server_gin/signutil"

	"github.com/gin-gonic/gin"
)

// Sign 校验请求的 HMAC-SHA256 签名。
func Sign(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		sign := ""
		if s, ok := params["sign"]; ok {
			sign = s.(string)
		}
		if sign == "" {
			sign = c.GetHeader("x-sign")
		}
		if sign == "" {
			sign = c.GetHeader("x-signature")
		}
		if sign == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "非法调用"})
			return
		}

		delete(params, "sign")
		payload := signutil.BuildParams(params)

		for _, secret := range cfg.App.Keys {
			if secret != "" && signutil.SHA256Hex(payload, secret) == sign {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": -1, "msg": "非法调用"})
	}
}

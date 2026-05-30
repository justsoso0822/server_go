package middleware

import (
	"net/http"

	"server_gin/config"
	signutil "server_gin/tools/sign"

	"github.com/gin-gonic/gin"
)

// Sign 校验请求的 HMAC-SHA256 签名。
func Sign(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := collectRequestParams(c)
		sign := requestSign(c, params)
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

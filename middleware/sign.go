package middleware

import (
	"net/http"

	"server_go/config"
	signutil "server_go/tools/sign"

	"github.com/gin-gonic/gin"
)

// 校验请求的 HMAC-SHA256 签名。
func Sign(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := collectRequestParams(c)
		sign := requestSign(c, params)
		// 存入 context 供下游中间件（ReplayGuard）复用，避免重复解析。
		c.Set("_params", params)
		c.Set("_sign", sign)
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

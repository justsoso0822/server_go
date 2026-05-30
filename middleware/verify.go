package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"server_gin/service"

	"github.com/gin-gonic/gin"
)

// Verify 校验 login_key，登录接口本身跳过。
func Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/user/login") {
			c.Next()
			return
		}

		uid := int64(0)
		if v := firstStr(c.Query("uid"), c.PostForm("uid")); v != "" {
			uid, _ = strconv.ParseInt(v, 10, 64)
		}

		result := service.VerifyLoginKey(c.Request.Context(), service.AuthInput{
			Uid:      uid,
			LoginKey: firstStr(c.Query("login_key"), c.PostForm("login_key")),
			Platform: firstStr(c.Query("platform"), c.PostForm("platform")),
			Version:  firstStr(c.Query("version"), c.PostForm("version")),
		})
		if result.Code != 0 {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": result.Code, "msg": result.Msg})
			return
		}

		c.Next()
	}
}

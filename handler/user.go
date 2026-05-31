package handler

import (
	"strconv"

	"server_go/service"
	"server_go/tools/autodb"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	ctx := c.Request.Context()
	uid, _ := strconv.ParseInt(firstParam(c, "uid"), 10, 64)
	loginKey := firstParam(c, "login_key")
	openid := firstParam(c, "openid")
	platform := firstParam(c, "platform")
	version := firstParam(c, "version")

	if uid == 0 || loginKey == "" || openid == "" || platform == "" || version == "" {
		fail(c, errParam)
		return
	}
	ctx = autodb.WithLogIdentity(ctx, uid, openid)
	c.Request = c.Request.WithContext(ctx)

	out, err := service.UserLogin(ctx, uid, loginKey, openid, platform, version)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, out)
}

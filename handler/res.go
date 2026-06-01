package handler

import (
	"server_go/service"

	"github.com/gin-gonic/gin"
)

func AddTili(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	if uid == 0 {
		fail(c, errParam)
		return
	}
	out, err := service.UpdateTili(ctx, uid, 50, "测试增加体力")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"res": out["res"], "__add_tili": out["add_value"]})
}

func AddGold(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	if uid == 0 {
		fail(c, errParam)
		return
	}
	out, err := service.UpdateGold(ctx, uid, 50, "测试增加金币")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"res": out["res"], "__add_gold": out["add_value"]})
}

func AddDiamond(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	if uid == 0 {
		fail(c, errParam)
		return
	}
	out, err := service.UpdateDiamond(ctx, uid, 50, "测试增加钻石")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"res": out["res"], "__add_diamond": out["add_value"]})
}

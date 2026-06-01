package handler

import (
	"strconv"

	"server_go/service"

	"github.com/gin-gonic/gin"
)

func GetBag(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	chapter, _ := strconv.Atoi(c.Param("chapter"))

	if uid == 0 {
		fail(c, errParam)
		return
	}

	out, err := service.GetUserBag(ctx, uid, chapter)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, out)
}

func GetBagTp(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	chapter, _ := strconv.Atoi(c.Param("chapter"))

	if uid == 0 {
		fail(c, errParam)
		return
	}

	out, err := service.GetUserBagTp(ctx, uid, chapter)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, out)
}

func GetGrid(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	chapter, _ := strconv.Atoi(c.Param("chapter"))

	if uid == 0 {
		fail(c, errParam)
		return
	}

	out, err := service.GetGrid(ctx, uid, chapter)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, out)
}

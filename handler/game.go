package handler

import (
	"strconv"
	"time"

	"server_go/service"

	"github.com/gin-gonic/gin"
)

func GameTime(c *gin.Context) {
	ok(c, gin.H{"now": time.Now().UnixMilli()})
}

func GameOnline(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	seconds, _ := strconv.ParseInt(firstParam(c, "seconds"), 10, 64)

	if uid == 0 || seconds < 0 {
		fail(c, errParam)
		return
	}

	if err := service.GameOnline(ctx, uid, seconds); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"now": time.Now().UnixMilli()})
}

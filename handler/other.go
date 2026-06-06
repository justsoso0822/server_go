package handler

import (
	"server_go/service"
	"server_go/tools/autodb"

	"github.com/gin-gonic/gin"
)

func ResVersion(c *gin.Context) {
	ctx := c.Request.Context()
	key := c.Param("key")
	if key == "" {
		fail(c, errParam)
		return
	}
	out, err := service.GetResVersion(ctx, key)
	if err != nil {
		fail(c, err.Error())
		return
	}
	rawOK(c, out)
}

func TestIndex(c *gin.Context) {
	ctx := c.Request.Context()
	var rows []map[string]interface{}
	autodb.DB(ctx).Raw(`SELECT u.uid, u.openid, log.time FROM user u
		LEFT JOIN log_login log ON u.uid = log.uid
		WHERE u.uid = ?
		ORDER BY log.time DESC`, 13081).Scan(&rows)
	ok(c, rows)
}

func TestDB(c *gin.Context) {
	ctx := c.Request.Context()
	uid := parseUID(firstParam(c, "uid"))
	if uid == 0 {
		fail(c, errParam)
		return
	}
	var row map[string]interface{}
	autodb.DB(ctx).Raw("SELECT * FROM user WHERE uid = ?", uid).Scan(&row)
	ok(c, row)
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const errParam = "参数错误"

func firstParam(c interface {
	Query(string) string
	PostForm(string) string
}, key string) string {
	if v := c.Query(key); v != "" {
		return v
	}
	return c.PostForm(key)
}

// parseUID 将 uid 字符串转为 int64，超过 int32 上限时钳位避免后续转换溢出。
func parseUID(v string) int64 {
	uid, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return uid
}

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

func fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": msg})
}

func rawOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

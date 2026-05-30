package handler

import (
	"net/http"

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

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

func fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": msg})
}

func rawOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

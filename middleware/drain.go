package middleware

import (
	"net/http"

	"server_go/state"

	"github.com/gin-gonic/gin"
)

// 在排水阶段拒绝新请求，并维护在途请求计数。
func DrainGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if state.IsRejecting() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"code": -1, "msg": "service is draining"})
			return
		}
		state.IncActiveRequests()
		defer state.DecActiveRequests()
		c.Next()
	}
}

package middleware

import (
	"net/http"

	"app/state"

	"github.com/gin-gonic/gin"
)

// 在排水阶段拒绝新请求，并维护在途请求计数。
func DrainGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		state.IncActiveRequests()
		defer state.DecActiveRequests()

		if state.IsRejecting() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"code": -1, "msg": "service is draining"})
			return
		}

		c.Next()
	}
}

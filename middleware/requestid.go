package middleware

import (
	"strings"

	"server_gin/tools/autodb"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-Id"
const maxRequestIDLen = 128

// RequestID propagates an inbound request id or creates one for application logs.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if requestID == "" || len(requestID) > maxRequestIDLen {
			requestID = uuid.NewString()
		}

		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Request.Header.Set(requestIDHeader, requestID)
		c.Request = c.Request.WithContext(autodb.WithRequestID(c.Request.Context(), requestID))
		c.Next()
	}
}

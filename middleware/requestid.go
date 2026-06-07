package middleware

import (
	"server_go/tools/autodb"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-Id"

// RequestID creates a server-side request id for application logs.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Request = c.Request.WithContext(autodb.WithRequestID(c.Request.Context(), requestID))
		c.Next()
	}
}

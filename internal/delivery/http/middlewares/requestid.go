package middlewares

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	requestID, ok := c.Get(RequestIDKey)
	if !ok {
		return ""
	}

	value, ok := requestID.(string)
	if !ok {
		return ""
	}

	return value
}

func newRequestID() string {
	buffer := make([]byte, 12)
	if _, err := rand.Read(buffer); err != nil {
		return "unknown"
	}

	return hex.EncodeToString(buffer)
}

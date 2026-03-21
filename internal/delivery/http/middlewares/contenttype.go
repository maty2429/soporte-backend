package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ContentTypeJSON validates that requests with body (POST, PUT, PATCH)
// have Content-Type: application/json header.
func ContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate methods that typically have a body
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
			c.Next()
			return
		}

		// Skip if no body
		if c.Request.ContentLength == 0 {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")

		// Check for application/json (with optional charset)
		if !isJSONContentType(contentType) {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
				"error": gin.H{
					"code":    "unsupported_media_type",
					"message": "Content-Type must be application/json",
				},
				"request_id": GetRequestID(c),
			})
			return
		}

		c.Next()
	}
}

func isJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	// Handle "application/json" or "application/json; charset=utf-8"
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	return mediaType == "application/json"
}

package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		// Get stack trace but limit depth to avoid exposing too much
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)
		stackTrace := string(stack[:length])

		// Log panic with sanitized info - don't log the actual panic value
		// as it might contain sensitive data
		log.Error("panic recovered",
			"panic_type", fmt.Sprintf("%T", recovered),
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"request_id", GetRequestID(c),
			"stack", stackTrace,
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "internal server error",
			},
			"request_id": GetRequestID(c),
		})
	})
}

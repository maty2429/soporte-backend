package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// QueryTimeout adds a timeout to the request context.
// This timeout is propagated to database queries via WithContext(ctx).
// If the handler returns without writing a response and the deadline was
// exceeded, the middleware writes a 504 Gateway Timeout.
//
// Note: this runs c.Next() synchronously (no goroutine) so the gin.Context
// is always accessed from a single goroutine, which is safe under -race.
func QueryTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout <= 0 {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if !c.Writer.Written() && ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": gin.H{
					"code":    "timeout",
					"message": "request timeout exceeded",
				},
				"request_id": GetRequestID(c),
			})
		}
	}
}

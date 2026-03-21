package middlewares

import (
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// safeQueryParams defines query parameters that are safe to log.
// Any parameter not in this list will be redacted.
var safeQueryParams = map[string]bool{
	"limit":  true,
	"offset": true,
	"page":   true,
	"size":   true,
	"sort":   true,
	"order":  true,
	"q":      true,
	"search": true,
	"filter": true,
	"id":     true,
	"ids":    true,
	"activa": true,
	"status": true,
	"type":   true,
	"from":   true,
	"to":     true,
}

func StructuredLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		path := c.Request.URL.Path

		c.Next()

		sanitizedQuery := sanitizeQueryParams(c.Request.URL.Query())
		if sanitizedQuery != "" {
			path = path + "?" + sanitizedQuery
		}

		status := c.Writer.Status()
		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency", time.Since(startedAt).String(),
			"client_ip", c.ClientIP(),
			"request_id", GetRequestID(c),
			"bytes_out", c.Writer.Size(),
		}

		switch {
		case status >= 500:
			log.Error("http request", attrs...)
		case status >= 400:
			log.Warn("http request", attrs...)
		default:
			log.Debug("http request", attrs...)
		}
	}
}

// sanitizeQueryParams removes sensitive query parameters from logs.
// Only parameters in safeQueryParams are included, others are redacted.
func sanitizeQueryParams(params url.Values) string {
	if len(params) == 0 {
		return ""
	}

	sanitized := make(url.Values)
	hasRedacted := false

	for key, values := range params {
		lowerKey := strings.ToLower(key)
		if safeQueryParams[lowerKey] {
			sanitized[key] = values
		} else {
			hasRedacted = true
		}
	}

	result := sanitized.Encode()
	if hasRedacted {
		if result != "" {
			result += "&"
		}
		result += "_redacted=true"
	}

	return result
}

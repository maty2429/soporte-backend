package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/config"
)

const (
	corsMaxAge         = "86400" // 24 h preflight cache
	corsMethods        = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsHeaders        = "Authorization, Content-Type, X-Request-ID"
	corsExposeHeaders  = "X-Request-ID"
)

// CORS handles cross-origin resource sharing.
//
// Behaviour:
//   - CORSAllowAll=true  → Access-Control-Allow-Origin: * (no credentials)
//   - CORSAllowAll=false → validate the incoming Origin against CORSAllowedOrigins;
//     if it matches, reflect that origin and add Vary: Origin.
//     Requests from unlisted origins receive no ACAO header (browser blocks them).
//   - Preflight (OPTIONS) → respond 204 and abort; no downstream handlers run.
func CORS(cfg config.SecurityConfig) gin.HandlerFunc {
	// Build a fast lookup set from the allowed-origins list.
	allowed := make(map[string]struct{}, len(cfg.CORSAllowedOrigins))
	for _, o := range cfg.CORSAllowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		switch {
		case cfg.CORSAllowAll:
			c.Header("Access-Control-Allow-Origin", "*")

		case origin != "":
			if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				// Vary tells caches that the response differs by origin.
				c.Header("Vary", "Origin")
			}
			// Unlisted origins: no ACAO header → browser enforces same-origin.
		}

		c.Header("Access-Control-Allow-Methods", corsMethods)
		c.Header("Access-Control-Allow-Headers", corsHeaders)
		c.Header("Access-Control-Expose-Headers", corsExposeHeaders)
		c.Header("Access-Control-Max-Age", corsMaxAge)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

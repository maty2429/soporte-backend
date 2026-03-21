package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"soporte/internal/config"
	"soporte/internal/delivery/http/handlers"
	docs "soporte/internal/delivery/http/handlers/docs"
)

func registerInfraRoutes(router *gin.Engine, v1 *gin.RouterGroup, h handlers.HealthHandler, d docs.DocsHandler, cfg config.Config) {
	router.GET("/health", h.Get)
	router.GET("/livez", h.Livez)
	router.GET("/readyz", h.Readyz)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if cfg.Docs.Enabled {
		router.GET("/openapi.json", d.OpenAPI)
		router.GET("/swagger/*any", d.SwaggerUI())
	}

	v1.GET("/health", h.Get)
	v1.GET("/livez", h.Livez)
	v1.GET("/readyz", h.Readyz)
}

// registerRootRoute registra el endpoint "/" con la lista de rutas disponibles.
// Debe llamarse al final de NewRouter para que router.Routes() esté completo.
func registerRootRoute(router *gin.Engine, cfg config.Config) {
	// Agrupa los métodos HTTP por path: { "/api/v1/solicitantes": ["GET", "POST"], ... }
	grouped := make(map[string][]string)
	order := make([]string, 0)
	for _, r := range router.Routes() {
		if _, seen := grouped[r.Path]; !seen {
			order = append(order, r.Path)
		}
		grouped[r.Path] = append(grouped[r.Path], r.Method)
	}

	type routeInfo struct {
		Path    string   `json:"path"`
		Methods []string `json:"methods"`
	}
	routes := make([]routeInfo, 0, len(order))
	for _, path := range order {
		routes = append(routes, routeInfo{Path: path, Methods: grouped[path]})
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": cfg.App.Name,
			"version": cfg.App.Version,
			"message": "API base lista para crecer",
			"routes":  routes,
		})
	})
}

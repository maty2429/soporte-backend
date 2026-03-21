//go:build !production

package docs

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed openapi.json
var docsFS embed.FS

type DocsHandler struct{}

func NewDocsHandler() DocsHandler {
	return DocsHandler{}
}

func (h DocsHandler) OpenAPI(c *gin.Context) {
	content, err := docsFS.ReadFile("openapi.json")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", content)
}

func (h DocsHandler) SwaggerUI() gin.HandlerFunc {
	swaggerHandler := gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/openapi.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
	))

	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		swaggerHandler(c)
	}
}

func (h DocsHandler) SwaggerRedirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}

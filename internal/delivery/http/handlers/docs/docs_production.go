//go:build production

package docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocsHandler struct{}

func NewDocsHandler() DocsHandler {
	return DocsHandler{}
}

func (h DocsHandler) OpenAPI(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

func (h DocsHandler) SwaggerUI() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	}
}

func (h DocsHandler) SwaggerRedirect(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/application/services"
	"soporte/internal/delivery/http/contracts"
)

type CatalogoFallaHandler struct {
	service *services.CatalogoFallaService
}

func NewCatalogoFallaHandler(service *services.CatalogoFallaService) CatalogoFallaHandler {
	return CatalogoFallaHandler{service: service}
}

func (h CatalogoFallaHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewCatalogoFallasResponse(items))
}

func (h CatalogoFallaHandler) Create(c *gin.Context) {
	request, ok := bindJSON[contracts.CreateCatalogoFallaRequest](c)
	if !ok {
		return
	}

	item, err := h.service.Create(c.Request.Context(), services.CreateCatalogoFallaCommand{
		CodigoFalla:          request.CodigoFalla,
		DescripcionFalla:     request.DescripcionFalla,
		Complejidad:          request.Complejidad,
		RequiereVisitaFisica: request.RequiereVisitaFisica,
		IDDepartamento:       request.IDDepartamento,
		Categoria:            request.Categoria,
		Subcategoria:         request.Subcategoria,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, fmt.Sprintf("/api/v1/catalogo-fallas/%d", item.ID), contracts.NewCatalogoFallaResponse(item))
}

func (h CatalogoFallaHandler) Update(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.UpdateCatalogoFallaRequest](c)
	if !ok {
		return
	}

	item, err := h.service.Update(c.Request.Context(), services.UpdateCatalogoFallaCommand{
		ID:                   id,
		CodigoFalla:          request.CodigoFalla,
		DescripcionFalla:     request.DescripcionFalla,
		Complejidad:          request.Complejidad,
		RequiereVisitaFisica: request.RequiereVisitaFisica,
		IDDepartamento:       request.IDDepartamento,
		Categoria:            request.Categoria,
		Subcategoria:         request.Subcategoria,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewCatalogoFallaResponse(item))
}

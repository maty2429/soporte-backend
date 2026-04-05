package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/application/services"
	"soporte/internal/delivery/http/contracts"
)

type ServicioHandler struct {
	service *services.ServicioService
}

func NewServicioHandler(service *services.ServicioService) ServicioHandler {
	return ServicioHandler{service: service}
}

func (h ServicioHandler) List(c *gin.Context) {
	query, ok := bindQuery[contracts.ListServiciosQuery](c)
	if !ok {
		return
	}

	result, err := h.service.List(c.Request.Context(), services.ListServiciosQuery{
		Edificio:  query.Edificio,
		Piso:      query.Piso,
		Servicios: query.Servicios,
		Search:    query.Search,
		Limit:     query.Limit,
		Offset:    query.Offset,
	})
	if err != nil {
		fail(c, err)
		return
	}

	list(c, contracts.NewServiciosResponse(result.Items), result.Total, result.Limit, result.Offset)
}

func (h ServicioHandler) Create(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateServicioRequest](c)
	if !ok {
		return
	}

	item, err := h.service.Create(c.Request.Context(), services.CreateServicioCommand{
		Edificio:                req.Edificio,
		Piso:                    req.Piso,
		Servicios:               req.Servicios,
		Ubicacion:               req.Ubicacion,
		Unidades:                req.Unidades,
		IDNivelPrioridadDefault: req.IDNivelPrioridadDefault,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, fmt.Sprintf("/api/v1/servicios/%d", item.ID), contracts.NewServicioResponse(item))
}

func (h ServicioHandler) Update(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	req, ok := bindJSON[contracts.UpdateServicioRequest](c)
	if !ok {
		return
	}

	item, err := h.service.Update(c.Request.Context(), services.UpdateServicioCommand{
		ID:                      id,
		Edificio:                req.Edificio,
		Piso:                    req.Piso,
		Servicios:               req.Servicios,
		Ubicacion:               req.Ubicacion,
		Unidades:                req.Unidades,
		IDNivelPrioridadDefault: req.IDNivelPrioridadDefault,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewServicioResponse(item))
}

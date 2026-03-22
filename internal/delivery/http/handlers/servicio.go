package handlers

import (
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
	request, ok := bindJSON[contracts.CreateServicioRequest](c)
	if !ok {
		return
	}

	item, err := h.service.Create(c.Request.Context(), services.CreateServicioCommand{
		Edificio:                request.Edificio,
		Piso:                    request.Piso,
		Servicios:               request.Servicios,
		Ubicacion:               request.Ubicacion,
		Unidades:                request.Unidades,
		IDNivelPrioridadDefault: request.IDNivelPrioridadDefault,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusCreated, contracts.NewServicioResponse(item))
}

func (h ServicioHandler) Update(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.UpdateServicioRequest](c)
	if !ok {
		return
	}

	item, err := h.service.Update(c.Request.Context(), services.UpdateServicioCommand{
		ID:                      id,
		Edificio:                request.Edificio,
		Piso:                    request.Piso,
		Servicios:               request.Servicios,
		Ubicacion:               request.Ubicacion,
		Unidades:                request.Unidades,
		IDNivelPrioridadDefault: request.IDNivelPrioridadDefault,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewServicioResponse(item))
}

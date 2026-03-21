package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/application/services"
	"soporte/internal/delivery/http/contracts"
)

type SolicitanteHandler struct {
	service *services.SolicitanteService
}

func NewSolicitanteHandler(service *services.SolicitanteService) SolicitanteHandler {
	return SolicitanteHandler{service: service}
}

func (h SolicitanteHandler) List(c *gin.Context) {
	query, ok := bindQuery[contracts.ListSolicitantesQuery](c)
	if !ok {
		return
	}

	result, err := h.service.List(c.Request.Context(), services.ListSolicitantesQuery{
		Limit:  query.Limit,
		Offset: query.Offset,
		Search: query.Search,
		Estado: query.Estado,
	})
	if err != nil {
		fail(c, err)
		return
	}

	list(c, contracts.NewSolicitantesResponse(result.Items), result.Total, result.Limit, result.Offset)
}

func (h SolicitanteHandler) Get(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	sol, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewSolicitanteResponse(sol))
}

func (h SolicitanteHandler) GetByRut(c *gin.Context) {
	sol, err := h.service.GetByRut(c.Request.Context(), c.Param("rut"))
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewSolicitanteResponse(sol))
}

func (h SolicitanteHandler) Create(c *gin.Context) {
	request, ok := bindJSON[contracts.CreateSolicitanteRequest](c)
	if !ok {
		return
	}

	sol, err := h.service.Create(c.Request.Context(), services.CreateSolicitanteCommand{
		IDServicio:     request.IDServicio,
		Correo:         request.Correo,
		Rut:            request.Rut,
		Dv:             request.Dv,
		NombreCompleto: request.NombreCompleto,
		Anexo:          request.Anexo,
		Estado:         request.Estado,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, fmt.Sprintf("/api/v1/solicitantes/%d", sol.ID), contracts.NewSolicitanteCreatedResponse(sol.ID))
}

func (h SolicitanteHandler) Update(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.UpdateSolicitanteRequest](c)
	if !ok {
		return
	}

	sol, err := h.service.Update(c.Request.Context(), services.UpdateSolicitanteCommand{
		ID:             id,
		IDServicio:     request.IDServicio,
		Correo:         request.Correo,
		Rut:            request.Rut,
		Dv:             request.Dv,
		NombreCompleto: request.NombreCompleto,
		Anexo:          request.Anexo,
		Estado:         request.Estado,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewSolicitanteResponse(sol))
}

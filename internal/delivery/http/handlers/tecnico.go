package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/application/services"
	"soporte/internal/delivery/http/contracts"
)

type TecnicoHandler struct {
	service *services.TecnicoService
}

func NewTecnicoHandler(service *services.TecnicoService) TecnicoHandler {
	return TecnicoHandler{service: service}
}

func (h TecnicoHandler) List(c *gin.Context) {
	query, ok := bindQuery[contracts.ListTecnicosQuery](c)
	if !ok {
		return
	}

	result, err := h.service.List(c.Request.Context(), services.ListTecnicosQuery{
		Limit:                 query.Limit,
		Offset:                query.Offset,
		Search:                query.Search,
		Estado:                query.Estado,
		IDTipoTecnico:         query.IDTipoTecnico,
		IDDepartamentoSoporte: query.IDDepartamentoSoporte,
	})
	if err != nil {
		fail(c, err)
		return
	}

	list(c, contracts.NewTecnicosResponse(result.Items), result.Total, result.Limit, result.Offset)
}

func (h TecnicoHandler) Get(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	t, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewTecnicoResponse(t))
}

func (h TecnicoHandler) GetByRut(c *gin.Context) {
	t, err := h.service.GetByRut(c.Request.Context(), c.Param("rut"))
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewTecnicoResponse(t))
}

func (h TecnicoHandler) Create(c *gin.Context) {
	request, ok := bindJSON[contracts.CreateTecnicoRequest](c)
	if !ok {
		return
	}

	t, err := h.service.Create(c.Request.Context(), services.CreateTecnicoCommand{
		Rut:                   request.Rut,
		Dv:                    request.Dv,
		NombreCompleto:        request.NombreCompleto,
		IDTipoTecnico:         request.IDTipoTecnico,
		IDDepartamentoSoporte: request.IDDepartamentoSoporte,
		IDTipoTurno:           request.IDTipoTurno,
		Estado:                request.Estado,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, fmt.Sprintf("/api/v1/tecnicos/%d", t.ID), contracts.NewTecnicoCreatedResponse(t.ID))
}

func (h TecnicoHandler) Update(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.UpdateTecnicoRequest](c)
	if !ok {
		return
	}

	t, err := h.service.Update(c.Request.Context(), services.UpdateTecnicoCommand{
		ID:                    id,
		Rut:                   request.Rut,
		Dv:                    request.Dv,
		NombreCompleto:        request.NombreCompleto,
		IDTipoTecnico:         request.IDTipoTecnico,
		IDDepartamentoSoporte: request.IDDepartamentoSoporte,
		IDTipoTurno:           request.IDTipoTurno,
		Estado:                request.Estado,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewTecnicoResponse(t))
}

// --- Configuración Horarios Turno ---

func (h TecnicoHandler) ListHorariosTurno(c *gin.Context) {
	items, err := h.service.ListHorariosTurno(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewHorariosTurnoResponse(items))
}

func (h TecnicoHandler) CreateHorarioTurno(c *gin.Context) {
	request, ok := bindJSON[contracts.CreateHorarioTurnoRequest](c)
	if !ok {
		return
	}

	horario, err := h.service.CreateHorarioTurno(c.Request.Context(), services.CreateHorarioTurnoCommand{
		IDTipoTurno: request.IDTipoTurno,
		DiaSemana:   request.DiaSemana,
		HoraInicio:  request.HoraInicio,
		HoraFin:     request.HoraFin,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, fmt.Sprintf("/api/v1/configuracion-horarios-turno/%d", horario.ID), contracts.NewHorarioTurnoResponse(horario))
}

func (h TecnicoHandler) UpdateHorarioTurno(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.UpdateHorarioTurnoRequest](c)
	if !ok {
		return
	}

	horario, err := h.service.UpdateHorarioTurno(c.Request.Context(), services.UpdateHorarioTurnoCommand{
		ID:          id,
		IDTipoTurno: request.IDTipoTurno,
		DiaSemana:   request.DiaSemana,
		HoraInicio:  request.HoraInicio,
		HoraFin:     request.HoraFin,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewHorarioTurnoResponse(horario))
}

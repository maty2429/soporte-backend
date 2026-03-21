package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/application/services"
	"soporte/internal/core/domain"
	"soporte/internal/delivery/http/contracts"
)

type CatalogoHandler struct {
	service *services.CatalogoService
}

func NewCatalogoHandler(service *services.CatalogoService) CatalogoHandler {
	return CatalogoHandler{service: service}
}

// ==================== tipos_ticket ====================

func (h CatalogoHandler) ListTiposTicket(c *gin.Context) {
	items, err := h.service.ListTiposTicket(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, items)
}

func (h CatalogoHandler) CreateTipoTicket(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateTipoTicketRequest](c)
	if !ok {
		return
	}
	item, err := h.service.CreateTipoTicket(c.Request.Context(), domain.TipoTicket{
		CodTipoTicket: req.CodTipoTicket,
		Descripcion:   req.Descripcion,
	})
	if err != nil {
		fail(c, err)
		return
	}
	created(c, fmt.Sprintf("/api/v1/tipos-ticket/%d", item.ID), item)
}

func (h CatalogoHandler) UpdateTipoTicket(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	req, ok := bindJSON[contracts.UpdateTipoTicketRequest](c)
	if !ok {
		return
	}
	item, err := h.service.UpdateTipoTicket(c.Request.Context(), id, req.CodTipoTicket, req.Descripcion)
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, item)
}

// ==================== niveles_prioridad ====================

func (h CatalogoHandler) ListNivelesPrioridad(c *gin.Context) {
	items, err := h.service.ListNivelesPrioridad(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, items)
}

func (h CatalogoHandler) CreateNivelPrioridad(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateDescripcionRequest](c)
	if !ok {
		return
	}
	item, err := h.service.CreateNivelPrioridad(c.Request.Context(), domain.NivelPrioridad{
		Descripcion: req.Descripcion,
	})
	if err != nil {
		fail(c, err)
		return
	}
	created(c, fmt.Sprintf("/api/v1/niveles-prioridad/%d", item.ID), item)
}

func (h CatalogoHandler) UpdateNivelPrioridad(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	req, ok := bindJSON[contracts.UpdateDescripcionRequest](c)
	if !ok {
		return
	}
	item, err := h.service.UpdateNivelPrioridad(c.Request.Context(), id, req.Descripcion)
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, item)
}

// ==================== tipos_tecnico ====================

func (h CatalogoHandler) ListTiposTecnico(c *gin.Context) {
	items, err := h.service.ListTiposTecnico(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, items)
}

func (h CatalogoHandler) CreateTipoTecnico(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateDescripcionRequest](c)
	if !ok {
		return
	}
	item, err := h.service.CreateTipoTecnico(c.Request.Context(), domain.TipoTecnico{
		Descripcion: req.Descripcion,
	})
	if err != nil {
		fail(c, err)
		return
	}
	created(c, fmt.Sprintf("/api/v1/tipos-tecnico/%d", item.ID), item)
}

func (h CatalogoHandler) UpdateTipoTecnico(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	req, ok := bindJSON[contracts.UpdateDescripcionRequest](c)
	if !ok {
		return
	}
	item, err := h.service.UpdateTipoTecnico(c.Request.Context(), id, req.Descripcion)
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, item)
}

// ==================== departamentos_soporte ====================

func (h CatalogoHandler) ListDepartamentosSoporte(c *gin.Context) {
	items, err := h.service.ListDepartamentosSoporte(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, items)
}

func (h CatalogoHandler) CreateDepartamentoSoporte(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateDepartamentoSoporteRequest](c)
	if !ok {
		return
	}
	item, err := h.service.CreateDepartamentoSoporte(c.Request.Context(), domain.DepartamentoSoporte{
		CodDepartamento: req.CodDepartamento,
		Descripcion:     req.Descripcion,
	})
	if err != nil {
		fail(c, err)
		return
	}
	created(c, fmt.Sprintf("/api/v1/departamentos-soporte/%d", item.ID), item)
}

func (h CatalogoHandler) UpdateDepartamentoSoporte(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	req, ok := bindJSON[contracts.UpdateDepartamentoSoporteRequest](c)
	if !ok {
		return
	}
	item, err := h.service.UpdateDepartamentoSoporte(c.Request.Context(), id, req.CodDepartamento, req.Descripcion)
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, item)
}

// ==================== motivos_pausa ====================

func (h CatalogoHandler) ListMotivosPausa(c *gin.Context) {
	items, err := h.service.ListMotivosPausa(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, items)
}

func (h CatalogoHandler) CreateMotivoPausa(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateMotivoPausaRequest](c)
	if !ok {
		return
	}
	item, err := h.service.CreateMotivoPausa(c.Request.Context(), domain.MotivoPausa{
		MotivoPausa:          req.MotivoPausa,
		RequiereAutorizacion: ptrOrDefault(req.RequiereAutorizacion, true),
	})
	if err != nil {
		fail(c, err)
		return
	}
	created(c, fmt.Sprintf("/api/v1/motivos-pausa/%d", item.ID), item)
}

func (h CatalogoHandler) UpdateMotivoPausa(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	req, ok := bindJSON[contracts.UpdateMotivoPausaRequest](c)
	if !ok {
		return
	}
	item, err := h.service.UpdateMotivoPausa(c.Request.Context(), id, req.MotivoPausa, req.RequiereAutorizacion)
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, item)
}

// ==================== tipos_turno ====================

func (h CatalogoHandler) ListTiposTurno(c *gin.Context) {
	items, err := h.service.ListTiposTurno(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, items)
}

func (h CatalogoHandler) CreateTipoTurno(c *gin.Context) {
	req, ok := bindJSON[contracts.CreateTipoTurnoRequest](c)
	if !ok {
		return
	}
	item, err := h.service.CreateTipoTurno(c.Request.Context(), domain.TipoTurno{
		Nombre:      req.Nombre,
		Descripcion: ptrOrStr(req.Descripcion),
		Estado:      ptrOrDefault(req.Estado, true),
	})
	if err != nil {
		fail(c, err)
		return
	}
	created(c, fmt.Sprintf("/api/v1/tipos-turno/%d", item.ID), item)
}

func (h CatalogoHandler) UpdateTipoTurno(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	req, ok := bindJSON[contracts.UpdateTipoTurnoRequest](c)
	if !ok {
		return
	}
	item, err := h.service.UpdateTipoTurno(c.Request.Context(), id, req.Nombre, req.Descripcion, req.Estado)
	if err != nil {
		fail(c, err)
		return
	}
	json(c, http.StatusOK, item)
}

// --- helpers ---

func ptrOrDefault[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

func ptrOrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

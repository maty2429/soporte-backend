package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"soporte/internal/application/services"
	"soporte/internal/delivery/http/contracts"
)

type TicketHandler struct {
	service *services.TicketService
}

func NewTicketHandler(service *services.TicketService) TicketHandler {
	return TicketHandler{service: service}
}

func (h TicketHandler) GetByID(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	ticket, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewTicketResponse(ticket))
}

func (h TicketHandler) ListTickets(c *gin.Context) {
	query, ok := bindQuery[contracts.ListTicketsQuery](c)
	if !ok {
		return
	}

	result, err := h.service.ListTickets(c.Request.Context(), services.ListTicketsQuery{
		CodEstadoTicket:       query.CodEstadoTicket,
		IDTecnicoAsignado:     query.IDTecnicoAsignado,
		RutTecnico:            query.RutTecnico,
		IDSolicitante:         query.IDSolicitante,
		IDDepartamentoSoporte: query.IDDepartamentoSoporte,
		Critico:               query.Critico,
		Limit:                 query.Limit,
		Offset:                query.Offset,
	})
	if err != nil {
		fail(c, err)
		return
	}

	list(c, contracts.NewTicketsResponse(result.Items), result.Total, result.Limit, result.Offset)
}

func (h TicketHandler) UpdateTicket(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.UpdateTicketRequest](c)
	if !ok {
		return
	}

	ticket, err := h.service.UpdateTicket(c.Request.Context(), services.UpdateTicketCommand{
		IDTicket:              id,
		DetalleFallaReportada: request.DetalleFallaReportada,
		UbicacionObs:          request.UbicacionObs,
		Critico:               request.Critico,
		IDTipoTicket:          request.IDTipoTicket,
		IDDepartamentoSoporte: request.IDDepartamentoSoporte,
		IDServicio:            request.IDServicio,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewTicketResponse(ticket))
}

func (h TicketHandler) ListBitacora(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	items, err := h.service.ListBitacora(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewBitacorasResponse(items))
}

func (h TicketHandler) GetByNroTicket(c *gin.Context) {
	detalle, err := h.service.GetByNroTicket(c.Request.Context(), c.Param("nro"))
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewTicketDetalleResponse(detalle))
}

func (h TicketHandler) Create(c *gin.Context) {
	request, ok := bindJSON[contracts.CreateTicketRequest](c)
	if !ok {
		return
	}

	ticket, err := h.service.Create(c.Request.Context(), services.CreateTicketCommand{
		IDSolicitante:         request.IDSolicitante,
		IDServicio:            request.IDServicio,
		IDTipoTicket:          request.IDTipoTicket,
		IDNivelPrioridad:      request.IDNivelPrioridad,
		IDDepartamentoSoporte: request.IDDepartamentoSoporte,
		Critico:               request.Critico,
		DetalleFallaReportada: request.DetalleFallaReportada,
		UbicacionObs:          request.UbicacionObs,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, "/api/v1/tickets/"+ticket.NroTicket, contracts.NewTicketCreatedResponse(ticket))
}

func (h TicketHandler) Assign(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.AssignTicketRequest](c)
	if !ok {
		return
	}

	ticket, err := h.service.Assign(c.Request.Context(), services.AssignTicketCommand{
		IDTicket:          id,
		IDTecnicoAsignado: request.IDTecnicoAsignado,
		IDCatalogoFalla:   request.IDCatalogoFalla,
		IDNivelPrioridad:  request.IDNivelPrioridad,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, contracts.NewAssignTicketResponse(ticket))
}

func (h TicketHandler) ListPausas(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	query, ok := bindQuery[contracts.ListPausasQuery](c)
	if !ok {
		return
	}

	result, err := h.service.ListPausas(c.Request.Context(), services.ListPausasQuery{
		IDTicket: id,
		Estado:   query.Estado,
		Limit:    query.Limit,
		Offset:   query.Offset,
	})
	if err != nil {
		fail(c, err)
		return
	}

	list(c, contracts.NewPausasDetalleResponse(result.Items), result.Total, result.Limit, result.Offset)
}

func (h TicketHandler) CreatePausa(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.CreatePausaRequest](c)
	if !ok {
		return
	}

	pausa, err := h.service.CreatePausa(c.Request.Context(), services.CreatePausaCommand{
		IDTicket:       id,
		IDTecnicoPausa: request.IDTecnicoPausa,
		IDMotivoPausa:  request.IDMotivoPausa,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusCreated, contracts.NewPausaResponse(pausa))
}

func (h TicketHandler) ResolverPausa(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.ResolverPausaRequest](c)
	if !ok {
		return
	}

	err := h.service.ResolverPausa(c.Request.Context(), services.ResolverPausaCommand{
		IDPausa:             id,
		EstadoPausa:         request.EstadoPausa,
		IDTecnicoAutorizado: request.IDTecnicoAutorizado,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, gin.H{"message": "pausa resuelta correctamente"})
}

func (h TicketHandler) ReanudarTicket(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.ReanudarTicketRequest](c)
	if !ok {
		return
	}

	err := h.service.ReanudarTicket(c.Request.Context(), services.ReanudarTicketCommand{
		IDTicket:       id,
		IDTecnicoPausa: request.IDTecnicoPausa,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, gin.H{"message": "ticket reanudado correctamente"})
}

func (h TicketHandler) Close(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.CloseTicketRequest](c)
	if !ok {
		return
	}

	err := h.service.Close(c.Request.Context(), services.CloseTicketCommand{
		IDTicket:      id,
		IDSolicitante: request.IDSolicitante,
		Nota:          request.Nota,
		Comentarios:   request.Comentarios,
		Observacion:   request.Observacion,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, gin.H{"message": "ticket cerrado correctamente"})
}

func (h TicketHandler) ChangeEstado(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.ChangeEstadoRequest](c)
	if !ok {
		return
	}

	err := h.service.ChangeEstado(c.Request.Context(), services.ChangeEstadoCommand{
		IDTicket:        id,
		CodEstadoTicket: request.CodEstadoTicket,
		RutResponsable:  request.RutResponsable,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, gin.H{"message": "estado actualizado correctamente"})
}

func (h TicketHandler) CreateTraspaso(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.CreateTraspasoRequest](c)
	if !ok {
		return
	}

	traspaso, err := h.service.CreateTraspaso(c.Request.Context(), services.CreateTraspasoCommand{
		IDTicket:         id,
		IDTecnicoOrigen:  request.IDTecnicoOrigen,
		IDTecnicoDestino: request.IDTecnicoDestino,
		Motivo:           request.Motivo,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusCreated, contracts.NewTraspasoResponse(traspaso))
}

func (h TicketHandler) ResolverTraspaso(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	query, ok := bindQuery[contracts.ResolverTraspasoQuery](c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.ResolverTraspasoRequest](c)
	if !ok {
		return
	}

	err := h.service.ResolverTraspaso(c.Request.Context(), services.ResolverTraspasoCommand{
		IDTraspaso:           id,
		IDTecnicoDestino:     query.IDTecnicoDestino,
		EstadoTraspaso:       request.EstadoTraspaso,
		ComentarioResolucion: request.ComentarioResolucion,
	})
	if err != nil {
		fail(c, err)
		return
	}

	json(c, http.StatusOK, gin.H{"message": "traspaso resuelto correctamente"})
}

func (h TicketHandler) ListTraspasos(c *gin.Context) {
	idTecnicoDestino, ok := getID(c)
	if !ok {
		return
	}

	query, ok := bindQuery[contracts.ListTraspasosQuery](c)
	if !ok {
		return
	}

	result, err := h.service.ListTraspasos(c.Request.Context(), services.ListTraspasosQuery{
		IDTecnicoDestino: idTecnicoDestino,
		Estado:           query.Estado,
		Limit:            query.Limit,
		Offset:           query.Offset,
	})
	if err != nil {
		fail(c, err)
		return
	}

	list(c, contracts.NewTraspasosResponse(result.Items), result.Total, result.Limit, result.Offset)
}

func (h TicketHandler) CreateBitacora(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	request, ok := bindJSON[contracts.CreateBitacoraRequest](c)
	if !ok {
		return
	}

	bitacora, err := h.service.CreateBitacora(c.Request.Context(), services.CreateBitacoraCommand{
		IDTicket:   id,
		RutAutor:   request.RutAutor,
		Comentario: request.Comentario,
	})
	if err != nil {
		fail(c, err)
		return
	}

	created(c, "", contracts.NewBitacoraResponse(bitacora))
}

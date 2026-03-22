package types

import "soporte/internal/core/domain"

type CreateTicketCommand struct {
	IDSolicitante         int
	IDServicio            int
	IDTipoTicket          int
	IDNivelPrioridad      int
	IDDepartamentoSoporte int
	Critico               bool
	DetalleFallaReportada string
	UbicacionObs          *string
}

type AssignTicketCommand struct {
	IDTicket          int
	IDTecnicoAsignado int
	IDCatalogoFalla   int
	IDNivelPrioridad  int
}

type CreateBitacoraCommand struct {
	IDTicket   int
	RutAutor   string
	Comentario string
}

type ChangeEstadoCommand struct {
	IDTicket        int
	CodEstadoTicket string
	RutResponsable  string
}

type CreatePausaCommand struct {
	IDTicket       int
	IDTecnicoPausa int
	IDMotivoPausa  int
}

type ResolverPausaCommand struct {
	IDPausa             int
	EstadoPausa         string // APROBADA o RECHAZADA
	IDTecnicoAutorizado int
}

type ReanudarTicketCommand struct {
	IDTicket       int
	IDTecnicoPausa int
}

type ListPausasQuery struct {
	IDTicket int
	Estado   string
	Limit    int
	Offset   int
}

type ListPausasResult struct {
	Items  []domain.TicketPausa
	Total  int64
	Limit  int
	Offset int
}

type CreateTraspasoCommand struct {
	IDTicket         int
	IDTecnicoOrigen  int
	IDTecnicoDestino int
	Motivo           string
}

type ResolverTraspasoCommand struct {
	IDTraspaso           int
	EstadoTraspaso       string // ACEPTADO o RECHAZADO
	ComentarioResolucion string
}

type ListTraspasosQuery struct {
	IDTicket int
	Estado   string
	Limit    int
	Offset   int
}

type ListTraspasosResult struct {
	Items  []domain.TicketTraspaso
	Total  int64
	Limit  int
	Offset int
}

type UpdateTicketCommand struct {
	IDTicket              int
	DetalleFallaReportada *string
	UbicacionObs          *string
	Critico               *bool
	IDTipoTicket          *int
	IDDepartamentoSoporte *int
	IDServicio            *int
}

type ListTicketsQuery struct {
	CodEstadoTicket       string
	IDTecnicoAsignado     int
	IDSolicitante         int
	IDDepartamentoSoporte int
	Critico               *bool
	Limit                 int
	Offset                int
}

type ListTicketsResult struct {
	Items  []domain.Ticket
	Total  int64
	Limit  int
	Offset int
}

type CloseTicketCommand struct {
	IDTicket      int
	IDSolicitante int
	Nota          int
	Comentarios   string
	Observacion   string
}

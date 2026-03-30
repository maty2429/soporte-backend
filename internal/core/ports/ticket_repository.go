package ports

import (
	"context"

	"soporte/internal/core/domain"
)

type ListPausasFilters struct {
	IDTicket int
	Estado   string // PENDIENTE, APROBADA, RECHAZADA — vacío = todas
	Limit    int
	Offset   int
}

type ListTraspasosFilters struct {
	IDTecnicoDestino int
	Estado           string // PENDIENTE, ACEPTADO, RECHAZADO — vacío = todos
	Limit            int
	Offset           int
}

type ListTicketsFilters struct {
	CodEstadoTicket       string
	IDTecnicoAsignado     int
	RutTecnico            string
	DVTecnico             string
	IDSolicitante         int
	IDDepartamentoSoporte int
	Critico               *bool
	Limit                 int
	Offset                int
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *domain.Ticket) error
	Update(ctx context.Context, ticket *domain.Ticket) error
	UpdateFields(ctx context.Context, ticket *domain.Ticket, fields map[string]any) error
	GetByID(ctx context.Context, id int) (domain.Ticket, error)
	ListTickets(ctx context.Context, filters ListTicketsFilters) ([]domain.Ticket, int64, error)
	GetByNroTicket(ctx context.Context, nro string) (domain.Ticket, error)
	NroTicketExists(ctx context.Context, nro string) (bool, error)
	GetEstadoTicketByCod(ctx context.Context, cod string) (domain.EstadoTicket, error)
	CreateTrazabilidad(ctx context.Context, t *domain.TrazabilidadTicket) error
	ListTrazabilidad(ctx context.Context, idTicket int) ([]domain.TrazabilidadTicket, error)
	CreateBitacora(ctx context.Context, b *domain.BitacoraTicket) error
	ListBitacora(ctx context.Context, idTicket int) ([]domain.BitacoraTicket, error)
	CreateValorizacion(ctx context.Context, v *domain.Valorizacion) error
	CreatePausa(ctx context.Context, p *domain.TicketPausa) error
	GetPausaByID(ctx context.Context, id int) (domain.TicketPausa, error)
	GetPausaActiva(ctx context.Context, idTicket int) (domain.TicketPausa, error)
	UpdatePausa(ctx context.Context, p *domain.TicketPausa) error
	ListPausas(ctx context.Context, filters ListPausasFilters) ([]domain.TicketPausa, int64, error)
	CreateTraspaso(ctx context.Context, t *domain.TicketTraspaso) error
	GetTraspasoByID(ctx context.Context, id int) (domain.TicketTraspaso, error)
	GetTraspasoPendiente(ctx context.Context, idTicket int) (domain.TicketTraspaso, error)
	UpdateTraspaso(ctx context.Context, t *domain.TicketTraspaso) error
	ListTraspasos(ctx context.Context, filters ListTraspasosFilters) ([]domain.TicketTraspaso, int64, error)
	RunInTx(ctx context.Context, fn func(txRepo TicketRepository) error) error
}

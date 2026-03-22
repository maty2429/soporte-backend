package repository

import (
	"context"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

var errUnavailable = domain.ServiceUnavailableError("database is not available", nil)

type unavailableSolicitanteRepository struct{}

func (r unavailableSolicitanteRepository) List(context.Context, ports.ListSolicitantesFilters) ([]domain.Solicitante, int64, error) {
	return nil, 0, errUnavailable
}

func (r unavailableSolicitanteRepository) GetByID(context.Context, int) (domain.Solicitante, error) {
	return domain.Solicitante{}, errUnavailable
}

func (r unavailableSolicitanteRepository) GetByRutDV(context.Context, string, string) (domain.Solicitante, error) {
	return domain.Solicitante{}, errUnavailable
}

func (r unavailableSolicitanteRepository) Create(context.Context, *domain.Solicitante) error {
	return errUnavailable
}

func (r unavailableSolicitanteRepository) Update(context.Context, *domain.Solicitante) error {
	return errUnavailable
}

type unavailableCatalogoRepository struct{}

func (r unavailableCatalogoRepository) ListTiposTicket(context.Context) ([]domain.TipoTicket, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoRepository) CreateTipoTicket(context.Context, *domain.TipoTicket) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) UpdateTipoTicket(context.Context, *domain.TipoTicket) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) GetTipoTicketByID(context.Context, int) (domain.TipoTicket, error) {
	return domain.TipoTicket{}, errUnavailable
}

func (r unavailableCatalogoRepository) ListNivelesPrioridad(context.Context) ([]domain.NivelPrioridad, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoRepository) CreateNivelPrioridad(context.Context, *domain.NivelPrioridad) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) UpdateNivelPrioridad(context.Context, *domain.NivelPrioridad) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) GetNivelPrioridadByID(context.Context, int) (domain.NivelPrioridad, error) {
	return domain.NivelPrioridad{}, errUnavailable
}

func (r unavailableCatalogoRepository) ListTiposTecnico(context.Context) ([]domain.TipoTecnico, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoRepository) CreateTipoTecnico(context.Context, *domain.TipoTecnico) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) UpdateTipoTecnico(context.Context, *domain.TipoTecnico) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) GetTipoTecnicoByID(context.Context, int) (domain.TipoTecnico, error) {
	return domain.TipoTecnico{}, errUnavailable
}

func (r unavailableCatalogoRepository) ListDepartamentosSoporte(context.Context) ([]domain.DepartamentoSoporte, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoRepository) CreateDepartamentoSoporte(context.Context, *domain.DepartamentoSoporte) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) UpdateDepartamentoSoporte(context.Context, *domain.DepartamentoSoporte) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) GetDepartamentoSoporteByID(context.Context, int) (domain.DepartamentoSoporte, error) {
	return domain.DepartamentoSoporte{}, errUnavailable
}

func (r unavailableCatalogoRepository) ListMotivosPausa(context.Context) ([]domain.MotivoPausa, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoRepository) CreateMotivoPausa(context.Context, *domain.MotivoPausa) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) UpdateMotivoPausa(context.Context, *domain.MotivoPausa) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) GetMotivoPausaByID(context.Context, int) (domain.MotivoPausa, error) {
	return domain.MotivoPausa{}, errUnavailable
}

func (r unavailableCatalogoRepository) ListTiposTurno(context.Context) ([]domain.TipoTurno, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoRepository) CreateTipoTurno(context.Context, *domain.TipoTurno) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) UpdateTipoTurno(context.Context, *domain.TipoTurno) error {
	return errUnavailable
}
func (r unavailableCatalogoRepository) GetTipoTurnoByID(context.Context, int) (domain.TipoTurno, error) {
	return domain.TipoTurno{}, errUnavailable
}

// --- servicio ---

type unavailableServicioRepository struct{}

func (r unavailableServicioRepository) List(context.Context, ports.ListServiciosFilters) ([]domain.Servicio, int64, error) {
	return nil, 0, errUnavailable
}
func (r unavailableServicioRepository) GetByID(context.Context, int) (domain.Servicio, error) {
	return domain.Servicio{}, errUnavailable
}
func (r unavailableServicioRepository) Create(context.Context, *domain.Servicio) error {
	return errUnavailable
}
func (r unavailableServicioRepository) Update(context.Context, *domain.Servicio) error {
	return errUnavailable
}

// --- catalogo falla ---

type unavailableCatalogoFallaRepository struct{}

func (r unavailableCatalogoFallaRepository) List(context.Context) ([]domain.CatalogoFalla, error) {
	return nil, errUnavailable
}
func (r unavailableCatalogoFallaRepository) GetByID(context.Context, int) (domain.CatalogoFalla, error) {
	return domain.CatalogoFalla{}, errUnavailable
}
func (r unavailableCatalogoFallaRepository) Create(context.Context, *domain.CatalogoFalla) error {
	return errUnavailable
}
func (r unavailableCatalogoFallaRepository) Update(context.Context, *domain.CatalogoFalla) error {
	return errUnavailable
}

// --- tecnico ---

type unavailableTecnicoRepository struct{}

func (r unavailableTecnicoRepository) List(context.Context, ports.ListTecnicosFilters) ([]domain.Tecnico, int64, error) {
	return nil, 0, errUnavailable
}
func (r unavailableTecnicoRepository) GetByID(context.Context, int) (domain.Tecnico, error) {
	return domain.Tecnico{}, errUnavailable
}
func (r unavailableTecnicoRepository) GetByRutDV(context.Context, string, string) (domain.Tecnico, error) {
	return domain.Tecnico{}, errUnavailable
}
func (r unavailableTecnicoRepository) Create(context.Context, *domain.Tecnico) error {
	return errUnavailable
}
func (r unavailableTecnicoRepository) Update(context.Context, *domain.Tecnico) error {
	return errUnavailable
}
func (r unavailableTecnicoRepository) ListHorariosTurno(context.Context) ([]domain.ConfiguracionHorarioTurno, error) {
	return nil, errUnavailable
}
func (r unavailableTecnicoRepository) CreateHorarioTurno(context.Context, *domain.ConfiguracionHorarioTurno) error {
	return errUnavailable
}
func (r unavailableTecnicoRepository) UpdateHorarioTurno(context.Context, *domain.ConfiguracionHorarioTurno) error {
	return errUnavailable
}
func (r unavailableTecnicoRepository) GetHorarioTurnoByID(context.Context, int) (domain.ConfiguracionHorarioTurno, error) {
	return domain.ConfiguracionHorarioTurno{}, errUnavailable
}

type unavailableTicketRepository struct{}

func (r unavailableTicketRepository) RunInTx(_ context.Context, _ func(ports.TicketRepository) error) error {
	return errUnavailable
}
func (r unavailableTicketRepository) Create(context.Context, *domain.Ticket) error {
	return errUnavailable
}
func (r unavailableTicketRepository) Update(context.Context, *domain.Ticket) error {
	return errUnavailable
}
func (r unavailableTicketRepository) UpdateFields(context.Context, *domain.Ticket, map[string]any) error {
	return errUnavailable
}
func (r unavailableTicketRepository) ListTickets(context.Context, ports.ListTicketsFilters) ([]domain.Ticket, int64, error) {
	return nil, 0, errUnavailable
}
func (r unavailableTicketRepository) GetByID(context.Context, int) (domain.Ticket, error) {
	return domain.Ticket{}, errUnavailable
}
func (r unavailableTicketRepository) GetByNroTicket(context.Context, string) (domain.Ticket, error) {
	return domain.Ticket{}, errUnavailable
}
func (r unavailableTicketRepository) NroTicketExists(context.Context, string) (bool, error) {
	return false, errUnavailable
}
func (r unavailableTicketRepository) GetEstadoTicketByCod(context.Context, string) (domain.EstadoTicket, error) {
	return domain.EstadoTicket{}, errUnavailable
}
func (r unavailableTicketRepository) CreateTrazabilidad(context.Context, *domain.TrazabilidadTicket) error {
	return errUnavailable
}
func (r unavailableTicketRepository) ListTrazabilidad(context.Context, int) ([]domain.TrazabilidadTicket, error) {
	return nil, errUnavailable
}
func (r unavailableTicketRepository) CreateBitacora(context.Context, *domain.BitacoraTicket) error {
	return errUnavailable
}
func (r unavailableTicketRepository) ListBitacora(context.Context, int) ([]domain.BitacoraTicket, error) {
	return nil, errUnavailable
}
func (r unavailableTicketRepository) CreateValorizacion(context.Context, *domain.Valorizacion) error {
	return errUnavailable
}
func (r unavailableTicketRepository) CreatePausa(context.Context, *domain.TicketPausa) error {
	return errUnavailable
}
func (r unavailableTicketRepository) GetPausaByID(context.Context, int) (domain.TicketPausa, error) {
	return domain.TicketPausa{}, errUnavailable
}
func (r unavailableTicketRepository) GetPausaActiva(context.Context, int) (domain.TicketPausa, error) {
	return domain.TicketPausa{}, errUnavailable
}
func (r unavailableTicketRepository) UpdatePausa(context.Context, *domain.TicketPausa) error {
	return errUnavailable
}
func (r unavailableTicketRepository) ListPausas(context.Context, ports.ListPausasFilters) ([]domain.TicketPausa, int64, error) {
	return nil, 0, errUnavailable
}
func (r unavailableTicketRepository) CreateTraspaso(context.Context, *domain.TicketTraspaso) error {
	return errUnavailable
}
func (r unavailableTicketRepository) GetTraspasoByID(context.Context, int) (domain.TicketTraspaso, error) {
	return domain.TicketTraspaso{}, errUnavailable
}
func (r unavailableTicketRepository) GetTraspasoPendiente(context.Context, int) (domain.TicketTraspaso, error) {
	return domain.TicketTraspaso{}, errUnavailable
}
func (r unavailableTicketRepository) UpdateTraspaso(context.Context, *domain.TicketTraspaso) error {
	return errUnavailable
}
func (r unavailableTicketRepository) ListTraspasos(context.Context, ports.ListTraspasosFilters) ([]domain.TicketTraspaso, int64, error) {
	return nil, 0, errUnavailable
}

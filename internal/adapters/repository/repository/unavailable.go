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

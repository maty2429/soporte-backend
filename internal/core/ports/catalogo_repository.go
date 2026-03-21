package ports

import (
	"context"

	"soporte/internal/core/domain"
)

type CatalogoRepository interface {
	// tipos_ticket
	ListTiposTicket(ctx context.Context) ([]domain.TipoTicket, error)
	CreateTipoTicket(ctx context.Context, item *domain.TipoTicket) error
	UpdateTipoTicket(ctx context.Context, item *domain.TipoTicket) error

	// niveles_prioridad
	ListNivelesPrioridad(ctx context.Context) ([]domain.NivelPrioridad, error)
	CreateNivelPrioridad(ctx context.Context, item *domain.NivelPrioridad) error
	UpdateNivelPrioridad(ctx context.Context, item *domain.NivelPrioridad) error

	// tipo_tecnico
	ListTiposTecnico(ctx context.Context) ([]domain.TipoTecnico, error)
	CreateTipoTecnico(ctx context.Context, item *domain.TipoTecnico) error
	UpdateTipoTecnico(ctx context.Context, item *domain.TipoTecnico) error

	// departamentos_soporte
	ListDepartamentosSoporte(ctx context.Context) ([]domain.DepartamentoSoporte, error)
	CreateDepartamentoSoporte(ctx context.Context, item *domain.DepartamentoSoporte) error
	UpdateDepartamentoSoporte(ctx context.Context, item *domain.DepartamentoSoporte) error

	// motivos_pausa
	ListMotivosPausa(ctx context.Context) ([]domain.MotivoPausa, error)
	CreateMotivoPausa(ctx context.Context, item *domain.MotivoPausa) error
	UpdateMotivoPausa(ctx context.Context, item *domain.MotivoPausa) error

	// tipos_turno
	ListTiposTurno(ctx context.Context) ([]domain.TipoTurno, error)
	CreateTipoTurno(ctx context.Context, item *domain.TipoTurno) error
	UpdateTipoTurno(ctx context.Context, item *domain.TipoTurno) error

	// GetByID (usado internamente por el service en Update)
	GetTipoTicketByID(ctx context.Context, id int) (domain.TipoTicket, error)
	GetNivelPrioridadByID(ctx context.Context, id int) (domain.NivelPrioridad, error)
	GetTipoTecnicoByID(ctx context.Context, id int) (domain.TipoTecnico, error)
	GetDepartamentoSoporteByID(ctx context.Context, id int) (domain.DepartamentoSoporte, error)
	GetMotivoPausaByID(ctx context.Context, id int) (domain.MotivoPausa, error)
	GetTipoTurnoByID(ctx context.Context, id int) (domain.TipoTurno, error)
}

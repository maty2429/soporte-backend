package ports

import (
	"context"

	"soporte/internal/core/domain"
)

type ListTecnicosFilters struct {
	IDTipoTecnico         int
	IDDepartamentoSoporte int
	Estado                *bool
	Search                string
	Limit                 int
	Offset                int
}

type TecnicoRepository interface {
	List(ctx context.Context, filters ListTecnicosFilters) ([]domain.Tecnico, int64, error)
	GetByID(ctx context.Context, id int) (domain.Tecnico, error)
	GetByRutDV(ctx context.Context, rut, dv string) (domain.Tecnico, error)
	Create(ctx context.Context, tecnico *domain.Tecnico) error
	Update(ctx context.Context, tecnico *domain.Tecnico) error
	ListHorariosTurno(ctx context.Context) ([]domain.ConfiguracionHorarioTurno, error)
	CreateHorarioTurno(ctx context.Context, h *domain.ConfiguracionHorarioTurno) error
	UpdateHorarioTurno(ctx context.Context, h *domain.ConfiguracionHorarioTurno) error
	GetHorarioTurnoByID(ctx context.Context, id int) (domain.ConfiguracionHorarioTurno, error)
}

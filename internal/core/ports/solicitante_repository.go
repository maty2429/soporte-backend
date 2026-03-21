package ports

//go:generate mockgen -destination=mocks/mock_solicitante_repository.go -package=mocks . SolicitanteRepository

import (
	"context"

	"soporte/internal/core/domain"
)

type ListSolicitantesFilters struct {
	Limit  int
	Offset int
	Search string
	Estado *bool
}

type SolicitanteRepository interface {
	List(ctx context.Context, filters ListSolicitantesFilters) ([]domain.Solicitante, int64, error)
	GetByID(ctx context.Context, id int) (domain.Solicitante, error)
	GetByRutDV(ctx context.Context, rut, dv string) (domain.Solicitante, error)
	Create(ctx context.Context, solicitante *domain.Solicitante) error
	Update(ctx context.Context, solicitante *domain.Solicitante) error
}

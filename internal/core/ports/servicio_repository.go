package ports

import (
	"context"

	"soporte/internal/core/domain"
)

type ListServiciosFilters struct {
	Edificio  string
	Piso      *int
	Servicios string
	Search    string // busca en unidades y ubicacion
	Limit     int
	Offset    int
}

type ServicioRepository interface {
	List(ctx context.Context, filters ListServiciosFilters) ([]domain.Servicio, int64, error)
	Create(ctx context.Context, item *domain.Servicio) error
	Update(ctx context.Context, item *domain.Servicio) error
	GetByID(ctx context.Context, id int) (domain.Servicio, error)
}

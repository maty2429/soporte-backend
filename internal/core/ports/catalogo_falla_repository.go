package ports

import (
	"context"

	"soporte/internal/core/domain"
)

type CatalogoFallaRepository interface {
	List(ctx context.Context) ([]domain.CatalogoFalla, error)
	Create(ctx context.Context, item *domain.CatalogoFalla) error
	Update(ctx context.Context, item *domain.CatalogoFalla) error
	GetByID(ctx context.Context, id int) (domain.CatalogoFalla, error)
}

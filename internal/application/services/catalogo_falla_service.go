package services

import (
	"context"
	"errors"
	"strings"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type CatalogoFallaService struct {
	repo ports.CatalogoFallaRepository
}

func NewCatalogoFallaService(repo ports.CatalogoFallaRepository) *CatalogoFallaService {
	return &CatalogoFallaService{repo: repo}
}

func (s *CatalogoFallaService) List(ctx context.Context) ([]domain.CatalogoFalla, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return nil, err
		}
		return nil, domain.InternalError("list catalogo fallas", err)
	}
	return items, nil
}

func (s *CatalogoFallaService) Create(ctx context.Context, cmd CreateCatalogoFallaCommand) (domain.CatalogoFalla, error) {
	codigoFalla := strings.TrimSpace(cmd.CodigoFalla)
	if codigoFalla == "" {
		return domain.CatalogoFalla{}, domain.ValidationError("codigo_falla is required", nil)
	}
	descripcion := strings.TrimSpace(cmd.DescripcionFalla)
	if descripcion == "" {
		return domain.CatalogoFalla{}, domain.ValidationError("descripcion_falla is required", nil)
	}
	if cmd.Complejidad < 1 || cmd.Complejidad > 10 {
		return domain.CatalogoFalla{}, domain.ValidationError("complejidad must be between 1 and 10", nil)
	}

	item := domain.CatalogoFalla{
		CodigoFalla:          codigoFalla,
		DescripcionFalla:     descripcion,
		Complejidad:          cmd.Complejidad,
		RequiereVisitaFisica: cmd.RequiereVisitaFisica,
		IDDepartamento:       cmd.IDDepartamento,
		Categoria:            strings.TrimSpace(cmd.Categoria),
		Subcategoria:         strings.TrimSpace(cmd.Subcategoria),
	}

	if err := s.repo.Create(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.CatalogoFalla{}, domain.ConflictError("codigo_falla already exists", err)
		}
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.CatalogoFalla{}, err
		}
		return domain.CatalogoFalla{}, domain.InternalError("create catalogo falla", err)
	}

	return item, nil
}

func (s *CatalogoFallaService) Update(ctx context.Context, cmd UpdateCatalogoFallaCommand) (domain.CatalogoFalla, error) {
	item, err := s.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.CatalogoFalla{}, err
		}
		return domain.CatalogoFalla{}, domain.InternalError("get catalogo falla", err)
	}

	if cmd.CodigoFalla != nil {
		item.CodigoFalla = strings.TrimSpace(*cmd.CodigoFalla)
	}
	if cmd.DescripcionFalla != nil {
		item.DescripcionFalla = strings.TrimSpace(*cmd.DescripcionFalla)
	}
	if cmd.Complejidad != nil {
		if *cmd.Complejidad < 1 || *cmd.Complejidad > 10 {
			return domain.CatalogoFalla{}, domain.ValidationError("complejidad must be between 1 and 10", nil)
		}
		item.Complejidad = *cmd.Complejidad
	}
	if cmd.RequiereVisitaFisica != nil {
		item.RequiereVisitaFisica = *cmd.RequiereVisitaFisica
	}
	if cmd.IDDepartamento != nil {
		item.IDDepartamento = cmd.IDDepartamento
	}
	if cmd.Categoria != nil {
		item.Categoria = strings.TrimSpace(*cmd.Categoria)
	}
	if cmd.Subcategoria != nil {
		item.Subcategoria = strings.TrimSpace(*cmd.Subcategoria)
	}

	if err := s.repo.Update(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.CatalogoFalla{}, domain.ConflictError("codigo_falla already exists", err)
		}
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.CatalogoFalla{}, err
		}
		return domain.CatalogoFalla{}, domain.InternalError("update catalogo falla", err)
	}

	return item, nil
}

package services

import (
	"context"
	"errors"
	"strings"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type ServicioService struct {
	repo ports.ServicioRepository
}

func NewServicioService(repo ports.ServicioRepository) *ServicioService {
	return &ServicioService{repo: repo}
}

func (s *ServicioService) List(ctx context.Context, q ListServiciosQuery) (ListServiciosResult, error) {
	limit, offset := normalizePagination(q.Limit, q.Offset)

	items, total, err := s.repo.List(ctx, ports.ListServiciosFilters{
		Edificio:  strings.TrimSpace(q.Edificio),
		Piso:      q.Piso,
		Servicios: strings.TrimSpace(q.Servicios),
		Search:    strings.TrimSpace(q.Search),
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return ListServiciosResult{}, err
		}
		return ListServiciosResult{}, domain.InternalError("list servicios", err)
	}

	return ListServiciosResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *ServicioService) Create(ctx context.Context, cmd CreateServicioCommand) (domain.Servicio, error) {
	item := domain.Servicio{
		Edificio:                strings.TrimSpace(cmd.Edificio),
		Piso:                    cmd.Piso,
		Servicios:               strings.TrimSpace(cmd.Servicios),
		Ubicacion:               strings.TrimSpace(cmd.Ubicacion),
		Unidades:                strings.TrimSpace(cmd.Unidades),
		IDNivelPrioridadDefault: cmd.IDNivelPrioridadDefault,
	}

	if err := s.repo.Create(ctx, &item); err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Servicio{}, err
		}
		return domain.Servicio{}, domain.InternalError("create servicio", err)
	}

	return item, nil
}

func (s *ServicioService) Update(ctx context.Context, cmd UpdateServicioCommand) (domain.Servicio, error) {
	item, err := s.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Servicio{}, err
		}
		return domain.Servicio{}, domain.InternalError("get servicio", err)
	}

	if cmd.Edificio != nil {
		item.Edificio = strings.TrimSpace(*cmd.Edificio)
	}
	if cmd.Piso != nil {
		item.Piso = *cmd.Piso
	}
	if cmd.Servicios != nil {
		item.Servicios = strings.TrimSpace(*cmd.Servicios)
	}
	if cmd.Ubicacion != nil {
		item.Ubicacion = strings.TrimSpace(*cmd.Ubicacion)
	}
	if cmd.Unidades != nil {
		item.Unidades = strings.TrimSpace(*cmd.Unidades)
	}
	if cmd.IDNivelPrioridadDefault != nil {
		item.IDNivelPrioridadDefault = cmd.IDNivelPrioridadDefault
	}

	if err := s.repo.Update(ctx, &item); err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Servicio{}, err
		}
		return domain.Servicio{}, domain.InternalError("update servicio", err)
	}

	return item, nil
}

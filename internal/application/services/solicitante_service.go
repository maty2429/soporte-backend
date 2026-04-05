package services

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type SolicitanteService struct {
	repo ports.SolicitanteRepository
}

func NewSolicitanteService(repo ports.SolicitanteRepository) *SolicitanteService {
	return &SolicitanteService{repo: repo}
}

func (s *SolicitanteService) List(ctx context.Context, query ListSolicitantesQuery) (ListSolicitantesResult, error) {
	limit, offset := normalizePagination(query.Limit, query.Offset)

	items, total, err := s.repo.List(ctx, ports.ListSolicitantesFilters{
		Limit:  limit,
		Offset: offset,
		Search: strings.TrimSpace(query.Search),
		Estado: query.Estado,
	})
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return ListSolicitantesResult{}, err
		}
		return ListSolicitantesResult{}, domain.InternalError("list solicitantes", err)
	}

	return ListSolicitantesResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *SolicitanteService) GetByID(ctx context.Context, id int) (domain.Solicitante, error) {
	sol, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Solicitante{}, domain.NotFoundError("solicitante", err)
		}
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Solicitante{}, err
		}
		return domain.Solicitante{}, domain.InternalError("get solicitante", err)
	}

	return sol, nil
}

func (s *SolicitanteService) GetByRut(ctx context.Context, raw string) (domain.Solicitante, error) {
	rut, dv, err := splitRutDV(raw)
	if err != nil {
		return domain.Solicitante{}, err
	}

	sol, err := s.repo.GetByRutDV(ctx, rut, dv)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Solicitante{}, domain.NotFoundError("solicitante", err)
		}
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Solicitante{}, err
		}
		return domain.Solicitante{}, domain.InternalError("get solicitante by rut", err)
	}

	return sol, nil
}

func (s *SolicitanteService) Create(ctx context.Context, command CreateSolicitanteCommand) (domain.Solicitante, error) {
	if strings.TrimSpace(command.Rut) == "" {
		return domain.Solicitante{}, domain.ValidationError("rut is required", nil)
	}
	if strings.TrimSpace(command.Dv) == "" {
		return domain.Solicitante{}, domain.ValidationError("dv is required", nil)
	}
	if strings.TrimSpace(command.NombreCompleto) == "" {
		return domain.Solicitante{}, domain.ValidationError("nombre_completo is required", nil)
	}

	estado := true
	if command.Estado != nil {
		estado = *command.Estado
	}

	sol := domain.Solicitante{
		IDServicio:     command.IDServicio,
		Correo:         strings.TrimSpace(strings.ToLower(command.Correo)),
		Rut:            strings.TrimSpace(command.Rut),
		Dv:             strings.TrimSpace(command.Dv),
		NombreCompleto: strings.TrimSpace(command.NombreCompleto),
		Anexo:          command.Anexo,
		Estado:         estado,
	}

	if err := s.repo.Create(ctx, &sol); err != nil {
		if isDuplicateKeyError(err) {
			return domain.Solicitante{}, domain.ConflictError("rut or correo already exists", err)
		}
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Solicitante{}, err
		}
		return domain.Solicitante{}, domain.InternalError("create solicitante", err)
	}

	return sol, nil
}

func (s *SolicitanteService) Update(ctx context.Context, command UpdateSolicitanteCommand) (domain.Solicitante, error) {
	sol, err := s.GetByID(ctx, command.ID)
	if err != nil {
		return domain.Solicitante{}, err
	}

	if command.IDServicio == nil && command.Correo == nil && command.Rut == nil &&
		command.Dv == nil && command.NombreCompleto == nil && command.Anexo == nil && command.Estado == nil {
		return domain.Solicitante{}, domain.ValidationError("at least one field must be provided", nil)
	}

	if command.IDServicio != nil {
		sol.IDServicio = command.IDServicio
	}
	if command.Correo != nil {
		sol.Correo = strings.TrimSpace(strings.ToLower(*command.Correo))
	}
	if command.Rut != nil {
		sol.Rut = strings.TrimSpace(*command.Rut)
	}
	if command.Dv != nil {
		sol.Dv = strings.TrimSpace(*command.Dv)
	}
	if command.NombreCompleto != nil {
		sol.NombreCompleto = strings.TrimSpace(*command.NombreCompleto)
	}
	if command.Anexo != nil {
		sol.Anexo = command.Anexo
	}
	if command.Estado != nil {
		sol.Estado = *command.Estado
	}

	if err := s.repo.Update(ctx, &sol); err != nil {
		if isDuplicateKeyError(err) {
			return domain.Solicitante{}, domain.ConflictError("rut or correo already exists", err)
		}
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Solicitante{}, err
		}
		return domain.Solicitante{}, domain.InternalError("update solicitante", err)
	}

	return sol, nil
}


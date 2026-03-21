package services

import (
	"context"
	"errors"
	"strings"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type CatalogoService struct {
	repo ports.CatalogoRepository
}

func NewCatalogoService(repo ports.CatalogoRepository) *CatalogoService {
	return &CatalogoService{repo: repo}
}

// ==================== tipos_ticket ====================

func (s *CatalogoService) ListTiposTicket(ctx context.Context) ([]domain.TipoTicket, error) {
	items, err := s.repo.ListTiposTicket(ctx)
	if err != nil {
		return nil, wrapServiceError("list tipo ticket", err)
	}
	return items, nil
}

func (s *CatalogoService) CreateTipoTicket(ctx context.Context, item domain.TipoTicket) (domain.TipoTicket, error) {
	item.CodTipoTicket = strings.TrimSpace(item.CodTipoTicket)
	item.Descripcion = strings.TrimSpace(item.Descripcion)

	if item.CodTipoTicket == "" {
		return domain.TipoTicket{}, domain.ValidationError("cod_tipo_ticket is required", nil)
	}
	if item.Descripcion == "" {
		return domain.TipoTicket{}, domain.ValidationError("descripcion is required", nil)
	}

	if err := s.repo.CreateTipoTicket(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.TipoTicket{}, domain.ConflictError("tipo ticket already exists", err)
		}
		return domain.TipoTicket{}, wrapServiceError("create tipo ticket", err)
	}
	return item, nil
}

func (s *CatalogoService) UpdateTipoTicket(ctx context.Context, id int, cod *string, desc *string) (domain.TipoTicket, error) {
	if cod == nil && desc == nil {
		return domain.TipoTicket{}, domain.ValidationError("at least one field must be provided", nil)
	}

	existing, err := s.repo.GetTipoTicketByID(ctx, id)
	if err != nil {
		return domain.TipoTicket{}, wrapServiceError("get tipo ticket", err)
	}

	if cod != nil {
		existing.CodTipoTicket = strings.TrimSpace(*cod)
	}
	if desc != nil {
		existing.Descripcion = strings.TrimSpace(*desc)
	}

	if err := s.repo.UpdateTipoTicket(ctx, &existing); err != nil {
		if isDuplicateKeyError(err) {
			return domain.TipoTicket{}, domain.ConflictError("tipo ticket already exists", err)
		}
		return domain.TipoTicket{}, wrapServiceError("update tipo ticket", err)
	}
	return existing, nil
}

// ==================== niveles_prioridad ====================

func (s *CatalogoService) ListNivelesPrioridad(ctx context.Context) ([]domain.NivelPrioridad, error) {
	items, err := s.repo.ListNivelesPrioridad(ctx)
	if err != nil {
		return nil, wrapServiceError("list nivel de prioridad", err)
	}
	return items, nil
}

func (s *CatalogoService) CreateNivelPrioridad(ctx context.Context, item domain.NivelPrioridad) (domain.NivelPrioridad, error) {
	item.Descripcion = strings.TrimSpace(item.Descripcion)

	if item.Descripcion == "" {
		return domain.NivelPrioridad{}, domain.ValidationError("descripcion is required", nil)
	}

	if err := s.repo.CreateNivelPrioridad(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.NivelPrioridad{}, domain.ConflictError("nivel de prioridad already exists", err)
		}
		return domain.NivelPrioridad{}, wrapServiceError("create nivel de prioridad", err)
	}
	return item, nil
}

func (s *CatalogoService) UpdateNivelPrioridad(ctx context.Context, id int, desc *string) (domain.NivelPrioridad, error) {
	if desc == nil {
		return domain.NivelPrioridad{}, domain.ValidationError("at least one field must be provided", nil)
	}

	existing, err := s.repo.GetNivelPrioridadByID(ctx, id)
	if err != nil {
		return domain.NivelPrioridad{}, wrapServiceError("get nivel de prioridad", err)
	}

	existing.Descripcion = strings.TrimSpace(*desc)

	if err := s.repo.UpdateNivelPrioridad(ctx, &existing); err != nil {
		if isDuplicateKeyError(err) {
			return domain.NivelPrioridad{}, domain.ConflictError("nivel de prioridad already exists", err)
		}
		return domain.NivelPrioridad{}, wrapServiceError("update nivel de prioridad", err)
	}
	return existing, nil
}

// ==================== tipo_tecnico ====================

func (s *CatalogoService) ListTiposTecnico(ctx context.Context) ([]domain.TipoTecnico, error) {
	items, err := s.repo.ListTiposTecnico(ctx)
	if err != nil {
		return nil, wrapServiceError("list tipo técnico", err)
	}
	return items, nil
}

func (s *CatalogoService) CreateTipoTecnico(ctx context.Context, item domain.TipoTecnico) (domain.TipoTecnico, error) {
	item.Descripcion = strings.TrimSpace(item.Descripcion)

	if item.Descripcion == "" {
		return domain.TipoTecnico{}, domain.ValidationError("descripcion is required", nil)
	}

	if err := s.repo.CreateTipoTecnico(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.TipoTecnico{}, domain.ConflictError("tipo técnico already exists", err)
		}
		return domain.TipoTecnico{}, wrapServiceError("create tipo técnico", err)
	}
	return item, nil
}

func (s *CatalogoService) UpdateTipoTecnico(ctx context.Context, id int, desc *string) (domain.TipoTecnico, error) {
	if desc == nil {
		return domain.TipoTecnico{}, domain.ValidationError("at least one field must be provided", nil)
	}

	existing, err := s.repo.GetTipoTecnicoByID(ctx, id)
	if err != nil {
		return domain.TipoTecnico{}, wrapServiceError("get tipo técnico", err)
	}

	existing.Descripcion = strings.TrimSpace(*desc)

	if err := s.repo.UpdateTipoTecnico(ctx, &existing); err != nil {
		if isDuplicateKeyError(err) {
			return domain.TipoTecnico{}, domain.ConflictError("tipo técnico already exists", err)
		}
		return domain.TipoTecnico{}, wrapServiceError("update tipo técnico", err)
	}
	return existing, nil
}

// ==================== departamentos_soporte ====================

func (s *CatalogoService) ListDepartamentosSoporte(ctx context.Context) ([]domain.DepartamentoSoporte, error) {
	items, err := s.repo.ListDepartamentosSoporte(ctx)
	if err != nil {
		return nil, wrapServiceError("list departamento de soporte", err)
	}
	return items, nil
}

func (s *CatalogoService) CreateDepartamentoSoporte(ctx context.Context, item domain.DepartamentoSoporte) (domain.DepartamentoSoporte, error) {
	item.CodDepartamento = strings.TrimSpace(item.CodDepartamento)
	item.Descripcion = strings.TrimSpace(item.Descripcion)

	if item.CodDepartamento == "" {
		return domain.DepartamentoSoporte{}, domain.ValidationError("cod_departamento is required", nil)
	}
	if item.Descripcion == "" {
		return domain.DepartamentoSoporte{}, domain.ValidationError("descripcion is required", nil)
	}

	if err := s.repo.CreateDepartamentoSoporte(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.DepartamentoSoporte{}, domain.ConflictError("departamento de soporte already exists", err)
		}
		return domain.DepartamentoSoporte{}, wrapServiceError("create departamento de soporte", err)
	}
	return item, nil
}

func (s *CatalogoService) UpdateDepartamentoSoporte(ctx context.Context, id int, cod *string, desc *string) (domain.DepartamentoSoporte, error) {
	if cod == nil && desc == nil {
		return domain.DepartamentoSoporte{}, domain.ValidationError("at least one field must be provided", nil)
	}

	existing, err := s.repo.GetDepartamentoSoporteByID(ctx, id)
	if err != nil {
		return domain.DepartamentoSoporte{}, wrapServiceError("get departamento de soporte", err)
	}

	if cod != nil {
		existing.CodDepartamento = strings.TrimSpace(*cod)
	}
	if desc != nil {
		existing.Descripcion = strings.TrimSpace(*desc)
	}

	if err := s.repo.UpdateDepartamentoSoporte(ctx, &existing); err != nil {
		if isDuplicateKeyError(err) {
			return domain.DepartamentoSoporte{}, domain.ConflictError("departamento de soporte already exists", err)
		}
		return domain.DepartamentoSoporte{}, wrapServiceError("update departamento de soporte", err)
	}
	return existing, nil
}

// ==================== motivos_pausa ====================

func (s *CatalogoService) ListMotivosPausa(ctx context.Context) ([]domain.MotivoPausa, error) {
	items, err := s.repo.ListMotivosPausa(ctx)
	if err != nil {
		return nil, wrapServiceError("list motivo de pausa", err)
	}
	return items, nil
}

func (s *CatalogoService) CreateMotivoPausa(ctx context.Context, item domain.MotivoPausa) (domain.MotivoPausa, error) {
	item.MotivoPausa = strings.TrimSpace(item.MotivoPausa)

	if item.MotivoPausa == "" {
		return domain.MotivoPausa{}, domain.ValidationError("motivo_pausa is required", nil)
	}
	if len(item.MotivoPausa) > 255 {
		return domain.MotivoPausa{}, domain.ValidationError("motivo_pausa exceeds max length of 255", nil)
	}

	if err := s.repo.CreateMotivoPausa(ctx, &item); err != nil {
		return domain.MotivoPausa{}, wrapServiceError("create motivo de pausa", err)
	}
	return item, nil
}

func (s *CatalogoService) UpdateMotivoPausa(ctx context.Context, id int, motivo *string, reqAuth *bool) (domain.MotivoPausa, error) {
	if motivo == nil && reqAuth == nil {
		return domain.MotivoPausa{}, domain.ValidationError("at least one field must be provided", nil)
	}

	existing, err := s.repo.GetMotivoPausaByID(ctx, id)
	if err != nil {
		return domain.MotivoPausa{}, wrapServiceError("get motivo de pausa", err)
	}

	if motivo != nil {
		trimmed := strings.TrimSpace(*motivo)
		if len(trimmed) > 255 {
			return domain.MotivoPausa{}, domain.ValidationError("motivo_pausa exceeds max length of 255", nil)
		}
		existing.MotivoPausa = trimmed
	}
	if reqAuth != nil {
		existing.RequiereAutorizacion = *reqAuth
	}

	if err := s.repo.UpdateMotivoPausa(ctx, &existing); err != nil {
		return domain.MotivoPausa{}, wrapServiceError("update motivo de pausa", err)
	}
	return existing, nil
}

// ==================== tipos_turno ====================

func (s *CatalogoService) ListTiposTurno(ctx context.Context) ([]domain.TipoTurno, error) {
	items, err := s.repo.ListTiposTurno(ctx)
	if err != nil {
		return nil, wrapServiceError("list tipo de turno", err)
	}
	return items, nil
}

func (s *CatalogoService) CreateTipoTurno(ctx context.Context, item domain.TipoTurno) (domain.TipoTurno, error) {
	item.Nombre = strings.TrimSpace(item.Nombre)
	item.Descripcion = strings.TrimSpace(item.Descripcion)

	if item.Nombre == "" {
		return domain.TipoTurno{}, domain.ValidationError("nombre is required", nil)
	}
	if len(item.Nombre) > 50 {
		return domain.TipoTurno{}, domain.ValidationError("nombre exceeds max length of 50", nil)
	}
	if len(item.Descripcion) > 255 {
		return domain.TipoTurno{}, domain.ValidationError("descripcion exceeds max length of 255", nil)
	}

	if err := s.repo.CreateTipoTurno(ctx, &item); err != nil {
		if isDuplicateKeyError(err) {
			return domain.TipoTurno{}, domain.ConflictError("tipo de turno already exists", err)
		}
		return domain.TipoTurno{}, wrapServiceError("create tipo de turno", err)
	}
	return item, nil
}

func (s *CatalogoService) UpdateTipoTurno(ctx context.Context, id int, nombre *string, desc *string, estado *bool) (domain.TipoTurno, error) {
	if nombre == nil && desc == nil && estado == nil {
		return domain.TipoTurno{}, domain.ValidationError("at least one field must be provided", nil)
	}

	existing, err := s.repo.GetTipoTurnoByID(ctx, id)
	if err != nil {
		return domain.TipoTurno{}, wrapServiceError("get tipo de turno", err)
	}

	if nombre != nil {
		trimmed := strings.TrimSpace(*nombre)
		if len(trimmed) > 50 {
			return domain.TipoTurno{}, domain.ValidationError("nombre exceeds max length of 50", nil)
		}
		existing.Nombre = trimmed
	}
	if desc != nil {
		trimmed := strings.TrimSpace(*desc)
		if len(trimmed) > 255 {
			return domain.TipoTurno{}, domain.ValidationError("descripcion exceeds max length of 255", nil)
		}
		existing.Descripcion = trimmed
	}
	if estado != nil {
		existing.Estado = *estado
	}

	if err := s.repo.UpdateTipoTurno(ctx, &existing); err != nil {
		if isDuplicateKeyError(err) {
			return domain.TipoTurno{}, domain.ConflictError("tipo de turno already exists", err)
		}
		return domain.TipoTurno{}, wrapServiceError("update tipo de turno", err)
	}
	return existing, nil
}

// --- helper ---

func wrapServiceError(op string, err error) error {
	var appErr *domain.Error
	if errors.As(err, &appErr) {
		return err
	}
	return domain.InternalError(op, err)
}

package services

import (
	"context"
	"strings"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type TecnicoService struct {
	repo ports.TecnicoRepository
}

func NewTecnicoService(repo ports.TecnicoRepository) *TecnicoService {
	return &TecnicoService{repo: repo}
}

func (s *TecnicoService) List(ctx context.Context, q ListTecnicosQuery) (ListTecnicosResult, error) {
	limit, offset := normalizePagination(q.Limit, q.Offset)

	items, total, err := s.repo.List(ctx, ports.ListTecnicosFilters{
		Limit:                 limit,
		Offset:                offset,
		Search:                strings.TrimSpace(q.Search),
		Estado:                q.Estado,
		IDTipoTecnico:         q.IDTipoTecnico,
		IDDepartamentoSoporte: q.IDDepartamentoSoporte,
	})
	if err != nil {
		return ListTecnicosResult{}, wrapServiceError("list tecnicos", err)
	}

	return ListTecnicosResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *TecnicoService) GetByID(ctx context.Context, id int) (domain.Tecnico, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Tecnico{}, wrapServiceError("get tecnico", err)
	}
	return t, nil
}

func (s *TecnicoService) GetByRut(ctx context.Context, raw string) (domain.Tecnico, error) {
	rut, dv, err := splitRutDV(raw)
	if err != nil {
		return domain.Tecnico{}, err
	}

	t, err := s.repo.GetByRutDV(ctx, rut, dv)
	if err != nil {
		return domain.Tecnico{}, wrapServiceError("get tecnico by rut", err)
	}
	return t, nil
}

func (s *TecnicoService) Create(ctx context.Context, cmd CreateTecnicoCommand) (domain.Tecnico, error) {
	if strings.TrimSpace(cmd.Rut) == "" {
		return domain.Tecnico{}, domain.ValidationError("rut is required", nil)
	}
	if strings.TrimSpace(cmd.Dv) == "" {
		return domain.Tecnico{}, domain.ValidationError("dv is required", nil)
	}
	if strings.TrimSpace(cmd.NombreCompleto) == "" {
		return domain.Tecnico{}, domain.ValidationError("nombre_completo is required", nil)
	}

	t := domain.Tecnico{
		Rut:                   strings.TrimSpace(cmd.Rut),
		Dv:                    strings.TrimSpace(cmd.Dv),
		NombreCompleto:        strings.TrimSpace(cmd.NombreCompleto),
		IDTipoTecnico:         cmd.IDTipoTecnico,
		IDDepartamentoSoporte: cmd.IDDepartamentoSoporte,
		Estado:                ptrOrDefaultBool(cmd.Estado, true),
	}

	if err := s.repo.Create(ctx, &t); err != nil {
		if isDuplicateKeyError(err) {
			return domain.Tecnico{}, domain.ConflictError("rut already exists", err)
		}
		return domain.Tecnico{}, wrapServiceError("create tecnico", err)
	}

	return t, nil
}

func (s *TecnicoService) Update(ctx context.Context, cmd UpdateTecnicoCommand) (domain.Tecnico, error) {
	t, err := s.GetByID(ctx, cmd.ID)
	if err != nil {
		return domain.Tecnico{}, err
	}

	if cmd.Rut == nil && cmd.Dv == nil && cmd.NombreCompleto == nil &&
		cmd.IDTipoTecnico == nil && cmd.IDDepartamentoSoporte == nil && cmd.Estado == nil {
		return domain.Tecnico{}, domain.ValidationError("at least one field must be provided", nil)
	}

	if cmd.Rut != nil {
		t.Rut = strings.TrimSpace(*cmd.Rut)
	}
	if cmd.Dv != nil {
		t.Dv = strings.TrimSpace(*cmd.Dv)
	}
	if cmd.NombreCompleto != nil {
		t.NombreCompleto = strings.TrimSpace(*cmd.NombreCompleto)
	}
	if cmd.IDTipoTecnico != nil {
		t.IDTipoTecnico = cmd.IDTipoTecnico
	}
	if cmd.IDDepartamentoSoporte != nil {
		t.IDDepartamentoSoporte = cmd.IDDepartamentoSoporte
	}
	if cmd.Estado != nil {
		t.Estado = *cmd.Estado
	}

	if err := s.repo.Update(ctx, &t); err != nil {
		if isDuplicateKeyError(err) {
			return domain.Tecnico{}, domain.ConflictError("rut already exists", err)
		}
		return domain.Tecnico{}, wrapServiceError("update tecnico", err)
	}

	return t, nil
}

// --- Configuración Horarios Turno ---

func (s *TecnicoService) ListHorariosTurno(ctx context.Context) ([]domain.ConfiguracionHorarioTurno, error) {
	items, err := s.repo.ListHorariosTurno(ctx)
	if err != nil {
		return nil, wrapServiceError("list horarios turno", err)
	}
	return items, nil
}

func (s *TecnicoService) CreateHorarioTurno(ctx context.Context, cmd CreateHorarioTurnoCommand) (domain.ConfiguracionHorarioTurno, error) {
	if cmd.IDTipoTurno <= 0 {
		return domain.ConfiguracionHorarioTurno{}, domain.ValidationError("id_tipo_turno is required", nil)
	}
	if cmd.DiaSemana < 0 || cmd.DiaSemana > 6 {
		return domain.ConfiguracionHorarioTurno{}, domain.ValidationError("dia_semana must be between 0 and 6", nil)
	}
	horaInicio := strings.TrimSpace(cmd.HoraInicio)
	if horaInicio == "" {
		return domain.ConfiguracionHorarioTurno{}, domain.ValidationError("hora_inicio is required", nil)
	}
	horaFin := strings.TrimSpace(cmd.HoraFin)
	if horaFin == "" {
		return domain.ConfiguracionHorarioTurno{}, domain.ValidationError("hora_fin is required", nil)
	}

	h := domain.ConfiguracionHorarioTurno{
		IDTipoTurno: cmd.IDTipoTurno,
		DiaSemana:   cmd.DiaSemana,
		HoraInicio:  horaInicio,
		HoraFin:     horaFin,
	}

	if err := s.repo.CreateHorarioTurno(ctx, &h); err != nil {
		return domain.ConfiguracionHorarioTurno{}, wrapServiceError("create horario turno", err)
	}

	return h, nil
}

func (s *TecnicoService) UpdateHorarioTurno(ctx context.Context, cmd UpdateHorarioTurnoCommand) (domain.ConfiguracionHorarioTurno, error) {
	h, err := s.repo.GetHorarioTurnoByID(ctx, cmd.ID)
	if err != nil {
		return domain.ConfiguracionHorarioTurno{}, wrapServiceError("get horario turno", err)
	}

	if cmd.IDTipoTurno == nil && cmd.DiaSemana == nil && cmd.HoraInicio == nil && cmd.HoraFin == nil {
		return domain.ConfiguracionHorarioTurno{}, domain.ValidationError("at least one field must be provided", nil)
	}

	if cmd.IDTipoTurno != nil {
		h.IDTipoTurno = *cmd.IDTipoTurno
	}
	if cmd.DiaSemana != nil {
		if *cmd.DiaSemana < 0 || *cmd.DiaSemana > 6 {
			return domain.ConfiguracionHorarioTurno{}, domain.ValidationError("dia_semana must be between 0 and 6", nil)
		}
		h.DiaSemana = *cmd.DiaSemana
	}
	if cmd.HoraInicio != nil {
		h.HoraInicio = strings.TrimSpace(*cmd.HoraInicio)
	}
	if cmd.HoraFin != nil {
		h.HoraFin = strings.TrimSpace(*cmd.HoraFin)
	}

	if err := s.repo.UpdateHorarioTurno(ctx, &h); err != nil {
		return domain.ConfiguracionHorarioTurno{}, wrapServiceError("update horario turno", err)
	}

	return h, nil
}

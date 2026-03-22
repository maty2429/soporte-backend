package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	dbmodels "soporte/internal/adapters/repository/models"
	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type TecnicoRepository struct {
	db *gorm.DB
}

func NewTecnicoRepository(db *gorm.DB) ports.TecnicoRepository {
	if db == nil {
		return unavailableTecnicoRepository{}
	}
	return &TecnicoRepository{db: db}
}

func (r *TecnicoRepository) List(ctx context.Context, f ports.ListTecnicosFilters) ([]domain.Tecnico, int64, error) {
	q := r.db.WithContext(ctx).Model(&dbmodels.Tecnico{})

	if f.Search != "" {
		like := "%" + strings.ToLower(f.Search) + "%"
		q = q.Where("LOWER(nombre_completo) LIKE ? OR rut LIKE ?", like, like)
	}
	if f.IDTipoTecnico > 0 {
		q = q.Where("id_tipo_tecnico = ?", f.IDTipoTecnico)
	}
	if f.IDDepartamentoSoporte > 0 {
		q = q.Where("id_departamento_soporte = ?", f.IDDepartamentoSoporte)
	}
	if f.Estado != nil {
		q = q.Where("estado = ?", *f.Estado)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, wrapListError("tecnicos", err)
	}

	var rows []dbmodels.Tecnico
	if err := q.Order("id ASC").Limit(f.Limit).Offset(f.Offset).Find(&rows).Error; err != nil {
		return nil, 0, wrapListError("tecnicos", err)
	}

	items := make([]domain.Tecnico, 0, len(rows))
	for _, row := range rows {
		items = append(items, toTecnicoDomain(row))
	}
	return items, total, nil
}

func (r *TecnicoRepository) GetByID(ctx context.Context, id int) (domain.Tecnico, error) {
	var row dbmodels.Tecnico
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return domain.Tecnico{}, wrapDBError("tecnico", err)
	}
	return toTecnicoDomain(row), nil
}

func (r *TecnicoRepository) GetByRutDV(ctx context.Context, rut, dv string) (domain.Tecnico, error) {
	var row dbmodels.Tecnico
	if err := r.db.WithContext(ctx).Where("rut = ? AND dv = ?", rut, dv).Take(&row).Error; err != nil {
		return domain.Tecnico{}, wrapDBError("tecnico", err)
	}
	return toTecnicoDomain(row), nil
}

func (r *TecnicoRepository) Create(ctx context.Context, t *domain.Tecnico) error {
	row := toTecnicoDB(*t)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("tecnico", err)
	}
	*t = toTecnicoDomain(row)
	return nil
}

func (r *TecnicoRepository) Update(ctx context.Context, t *domain.Tecnico) error {
	row := toTecnicoDB(*t)
	err := r.db.WithContext(ctx).
		Model(&row).
		Select("Rut", "Dv", "NombreCompleto", "IDTipoTecnico", "IDDepartamentoSoporte", "IDTipoTurno", "Estado").
		Updates(&row).Error
	return wrapUpdateError("tecnico", err)
}

func toTecnicoDomain(row dbmodels.Tecnico) domain.Tecnico {
	var estado bool
	if row.Estado != nil {
		estado = *row.Estado
	}
	return domain.Tecnico{
		ID:                    row.ID,
		Rut:                   row.Rut,
		Dv:                    row.Dv,
		NombreCompleto:        row.NombreCompleto,
		IDTipoTecnico:         row.IDTipoTecnico,
		IDDepartamentoSoporte: row.IDDepartamentoSoporte,
		IDTipoTurno:           row.IDTipoTurno,
		Estado:                estado,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}

// --- configuracion horarios turno ---

func (r *TecnicoRepository) ListHorariosTurno(ctx context.Context) ([]domain.ConfiguracionHorarioTurno, error) {
	var rows []dbmodels.ConfiguracionHorarioTurno
	if err := r.db.WithContext(ctx).Order("id_tipo_turno ASC, dia_semana ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("configuracion_horarios_turno", err)
	}
	items := make([]domain.ConfiguracionHorarioTurno, 0, len(rows))
	for _, row := range rows {
		items = append(items, toHorarioTurnoDomain(row))
	}
	return items, nil
}

func (r *TecnicoRepository) CreateHorarioTurno(ctx context.Context, h *domain.ConfiguracionHorarioTurno) error {
	row := dbmodels.ConfiguracionHorarioTurno{
		IDTipoTurno: h.IDTipoTurno,
		DiaSemana:   h.DiaSemana,
		HoraInicio:  h.HoraInicio,
		HoraFin:     h.HoraFin,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("configuracion_horarios_turno", err)
	}
	h.ID = row.ID
	return nil
}

func (r *TecnicoRepository) GetHorarioTurnoByID(ctx context.Context, id int) (domain.ConfiguracionHorarioTurno, error) {
	var row dbmodels.ConfiguracionHorarioTurno
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return domain.ConfiguracionHorarioTurno{}, wrapDBError("configuracion_horarios_turno", err)
	}
	return toHorarioTurnoDomain(row), nil
}

func (r *TecnicoRepository) UpdateHorarioTurno(ctx context.Context, h *domain.ConfiguracionHorarioTurno) error {
	row := dbmodels.ConfiguracionHorarioTurno{
		ID:          h.ID,
		IDTipoTurno: h.IDTipoTurno,
		DiaSemana:   h.DiaSemana,
		HoraInicio:  h.HoraInicio,
		HoraFin:     h.HoraFin,
	}
	err := r.db.WithContext(ctx).
		Model(&row).
		Select("IDTipoTurno", "DiaSemana", "HoraInicio", "HoraFin").
		Updates(&row).Error
	return wrapUpdateError("configuracion_horarios_turno", err)
}

func toHorarioTurnoDomain(row dbmodels.ConfiguracionHorarioTurno) domain.ConfiguracionHorarioTurno {
	return domain.ConfiguracionHorarioTurno{
		ID:          row.ID,
		IDTipoTurno: row.IDTipoTurno,
		DiaSemana:   row.DiaSemana,
		HoraInicio:  row.HoraInicio,
		HoraFin:     row.HoraFin,
	}
}

func toTecnicoDB(t domain.Tecnico) dbmodels.Tecnico {
	return dbmodels.Tecnico{
		ID:                    t.ID,
		Rut:                   t.Rut,
		Dv:                    t.Dv,
		NombreCompleto:        t.NombreCompleto,
		IDTipoTecnico:         t.IDTipoTecnico,
		IDDepartamentoSoporte: t.IDDepartamentoSoporte,
		IDTipoTurno:           t.IDTipoTurno,
		Estado:                &t.Estado,
	}
}

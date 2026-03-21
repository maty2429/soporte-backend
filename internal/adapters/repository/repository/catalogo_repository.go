package repository

import (
	"context"

	"gorm.io/gorm"

	"soporte/internal/adapters/repository/models"
	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type CatalogoRepository struct {
	db *gorm.DB
}

func NewCatalogoRepository(db *gorm.DB) ports.CatalogoRepository {
	if db == nil {
		return unavailableCatalogoRepository{}
	}
	return &CatalogoRepository{db: db}
}

// ==================== tipos_ticket ====================

func (r *CatalogoRepository) ListTiposTicket(ctx context.Context) ([]domain.TipoTicket, error) {
	var rows []models.TipoTicket
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("tipo ticket", err)
	}
	items := make([]domain.TipoTicket, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.TipoTicket{
			ID:            row.ID,
			CodTipoTicket: row.CodTipoTicket,
			Descripcion:   row.Descripcion,
		})
	}
	return items, nil
}

func (r *CatalogoRepository) GetTipoTicketByID(ctx context.Context, id int) (domain.TipoTicket, error) {
	var row models.TipoTicket
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.TipoTicket{}, wrapDBError("tipo ticket", err)
	}
	return domain.TipoTicket{
		ID:            row.ID,
		CodTipoTicket: row.CodTipoTicket,
		Descripcion:   row.Descripcion,
	}, nil
}

func (r *CatalogoRepository) CreateTipoTicket(ctx context.Context, item *domain.TipoTicket) error {
	row := models.TipoTicket{
		CodTipoTicket: item.CodTipoTicket,
		Descripcion:   item.Descripcion,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("tipo ticket", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoRepository) UpdateTipoTicket(ctx context.Context, item *domain.TipoTicket) error {
	updates := map[string]any{}
	if item.CodTipoTicket != "" {
		updates["cod_tipo_ticket"] = item.CodTipoTicket
	}
	if item.Descripcion != "" {
		updates["descripcion"] = item.Descripcion
	}
	err := r.db.WithContext(ctx).Model(&models.TipoTicket{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("tipo ticket", err)
}

// ==================== niveles_prioridad ====================

func (r *CatalogoRepository) ListNivelesPrioridad(ctx context.Context) ([]domain.NivelPrioridad, error) {
	var rows []models.NivelPrioridad
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("nivel de prioridad", err)
	}
	items := make([]domain.NivelPrioridad, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.NivelPrioridad{
			ID:          row.ID,
			Descripcion: row.Descripcion,
		})
	}
	return items, nil
}

func (r *CatalogoRepository) GetNivelPrioridadByID(ctx context.Context, id int) (domain.NivelPrioridad, error) {
	var row models.NivelPrioridad
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.NivelPrioridad{}, wrapDBError("nivel de prioridad", err)
	}
	return domain.NivelPrioridad{
		ID:          row.ID,
		Descripcion: row.Descripcion,
	}, nil
}

func (r *CatalogoRepository) CreateNivelPrioridad(ctx context.Context, item *domain.NivelPrioridad) error {
	row := models.NivelPrioridad{
		Descripcion: item.Descripcion,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("nivel de prioridad", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoRepository) UpdateNivelPrioridad(ctx context.Context, item *domain.NivelPrioridad) error {
	updates := map[string]any{}
	if item.Descripcion != "" {
		updates["descripcion"] = item.Descripcion
	}
	err := r.db.WithContext(ctx).Model(&models.NivelPrioridad{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("nivel de prioridad", err)
}

// ==================== tipo_tecnico ====================

func (r *CatalogoRepository) ListTiposTecnico(ctx context.Context) ([]domain.TipoTecnico, error) {
	var rows []models.TipoTecnico
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("tipo técnico", err)
	}
	items := make([]domain.TipoTecnico, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.TipoTecnico{
			ID:          row.ID,
			Descripcion: row.Descripcion,
		})
	}
	return items, nil
}

func (r *CatalogoRepository) GetTipoTecnicoByID(ctx context.Context, id int) (domain.TipoTecnico, error) {
	var row models.TipoTecnico
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.TipoTecnico{}, wrapDBError("tipo técnico", err)
	}
	return domain.TipoTecnico{
		ID:          row.ID,
		Descripcion: row.Descripcion,
	}, nil
}

func (r *CatalogoRepository) CreateTipoTecnico(ctx context.Context, item *domain.TipoTecnico) error {
	row := models.TipoTecnico{
		Descripcion: item.Descripcion,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("tipo técnico", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoRepository) UpdateTipoTecnico(ctx context.Context, item *domain.TipoTecnico) error {
	updates := map[string]any{}
	if item.Descripcion != "" {
		updates["descripcion"] = item.Descripcion
	}
	err := r.db.WithContext(ctx).Model(&models.TipoTecnico{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("tipo técnico", err)
}

// ==================== departamentos_soporte ====================

func (r *CatalogoRepository) ListDepartamentosSoporte(ctx context.Context) ([]domain.DepartamentoSoporte, error) {
	var rows []models.DepartamentoSoporte
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("departamento de soporte", err)
	}
	items := make([]domain.DepartamentoSoporte, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.DepartamentoSoporte{
			ID:              row.ID,
			CodDepartamento: row.CodDepartamento,
			Descripcion:     row.Descripcion,
		})
	}
	return items, nil
}

func (r *CatalogoRepository) GetDepartamentoSoporteByID(ctx context.Context, id int) (domain.DepartamentoSoporte, error) {
	var row models.DepartamentoSoporte
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.DepartamentoSoporte{}, wrapDBError("departamento de soporte", err)
	}
	return domain.DepartamentoSoporte{
		ID:              row.ID,
		CodDepartamento: row.CodDepartamento,
		Descripcion:     row.Descripcion,
	}, nil
}

func (r *CatalogoRepository) CreateDepartamentoSoporte(ctx context.Context, item *domain.DepartamentoSoporte) error {
	row := models.DepartamentoSoporte{
		CodDepartamento: item.CodDepartamento,
		Descripcion:     item.Descripcion,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("departamento de soporte", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoRepository) UpdateDepartamentoSoporte(ctx context.Context, item *domain.DepartamentoSoporte) error {
	updates := map[string]any{}
	if item.CodDepartamento != "" {
		updates["cod_departamento"] = item.CodDepartamento
	}
	if item.Descripcion != "" {
		updates["descripcion"] = item.Descripcion
	}
	err := r.db.WithContext(ctx).Model(&models.DepartamentoSoporte{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("departamento de soporte", err)
}

// ==================== motivos_pausa ====================

func (r *CatalogoRepository) ListMotivosPausa(ctx context.Context) ([]domain.MotivoPausa, error) {
	var rows []models.MotivoPausa
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("motivo de pausa", err)
	}
	items := make([]domain.MotivoPausa, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.MotivoPausa{
			ID:                   row.ID,
			MotivoPausa:          row.MotivoPausa,
			RequiereAutorizacion: derefBool(row.RequiereAutorizacion),
		})
	}
	return items, nil
}

func (r *CatalogoRepository) GetMotivoPausaByID(ctx context.Context, id int) (domain.MotivoPausa, error) {
	var row models.MotivoPausa
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.MotivoPausa{}, wrapDBError("motivo de pausa", err)
	}
	return domain.MotivoPausa{
		ID:                   row.ID,
		MotivoPausa:          row.MotivoPausa,
		RequiereAutorizacion: derefBool(row.RequiereAutorizacion),
	}, nil
}

func (r *CatalogoRepository) CreateMotivoPausa(ctx context.Context, item *domain.MotivoPausa) error {
	row := models.MotivoPausa{
		MotivoPausa:          item.MotivoPausa,
		RequiereAutorizacion: &item.RequiereAutorizacion,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("motivo de pausa", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoRepository) UpdateMotivoPausa(ctx context.Context, item *domain.MotivoPausa) error {
	updates := map[string]any{
		"requiere_autorizacion": item.RequiereAutorizacion,
	}
	if item.MotivoPausa != "" {
		updates["motivo_pausa"] = item.MotivoPausa
	}
	err := r.db.WithContext(ctx).Model(&models.MotivoPausa{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("motivo de pausa", err)
}

// ==================== tipos_turno ====================

func (r *CatalogoRepository) ListTiposTurno(ctx context.Context) ([]domain.TipoTurno, error) {
	var rows []models.TipoTurno
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("tipo de turno", err)
	}
	items := make([]domain.TipoTurno, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.TipoTurno{
			ID:          row.ID,
			Nombre:      row.Nombre,
			Descripcion: row.Descripcion,
			Estado:      derefBool(row.Estado),
		})
	}
	return items, nil
}

func (r *CatalogoRepository) GetTipoTurnoByID(ctx context.Context, id int) (domain.TipoTurno, error) {
	var row models.TipoTurno
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.TipoTurno{}, wrapDBError("tipo de turno", err)
	}
	return domain.TipoTurno{
		ID:          row.ID,
		Nombre:      row.Nombre,
		Descripcion: row.Descripcion,
		Estado:      derefBool(row.Estado),
	}, nil
}

func (r *CatalogoRepository) CreateTipoTurno(ctx context.Context, item *domain.TipoTurno) error {
	row := models.TipoTurno{
		Nombre:      item.Nombre,
		Descripcion: item.Descripcion,
		Estado:      &item.Estado,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("tipo de turno", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoRepository) UpdateTipoTurno(ctx context.Context, item *domain.TipoTurno) error {
	updates := map[string]any{
		"estado": item.Estado,
	}
	if item.Nombre != "" {
		updates["nombre"] = item.Nombre
	}
	if item.Descripcion != "" {
		updates["descripcion"] = item.Descripcion
	}
	err := r.db.WithContext(ctx).Model(&models.TipoTurno{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("tipo de turno", err)
}

// --- helpers ---

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

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

// listCatalogo is a generic helper for ordering and scanning catalog tables.
func listCatalogo[M any, D any](db *gorm.DB, ctx context.Context, resource string, mapFn func(M) D) ([]D, error) {
	var rows []M
	if err := db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError(resource, err)
	}
	items := make([]D, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapFn(row))
	}
	return items, nil
}

// getCatalogoByID is a generic helper for fetching a catalog row by primary key.
func getCatalogoByID[M any, D any](db *gorm.DB, ctx context.Context, id int, resource string, mapFn func(M) D) (D, error) {
	var row M
	if err := db.WithContext(ctx).First(&row, id).Error; err != nil {
		var zero D
		return zero, wrapDBError(resource, err)
	}
	return mapFn(row), nil
}

// ==================== tipos_ticket ====================

func (r *CatalogoRepository) ListTiposTicket(ctx context.Context) ([]domain.TipoTicket, error) {
	return listCatalogo(r.db, ctx, "tipo ticket", func(row models.TipoTicket) domain.TipoTicket {
		return domain.TipoTicket{ID: row.ID, CodTipoTicket: row.CodTipoTicket, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) GetTipoTicketByID(ctx context.Context, id int) (domain.TipoTicket, error) {
	return getCatalogoByID(r.db, ctx, id, "tipo ticket", func(row models.TipoTicket) domain.TipoTicket {
		return domain.TipoTicket{ID: row.ID, CodTipoTicket: row.CodTipoTicket, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) CreateTipoTicket(ctx context.Context, item *domain.TipoTicket) error {
	row := models.TipoTicket{CodTipoTicket: item.CodTipoTicket, Descripcion: item.Descripcion}
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
	return wrapUpdateError("tipo ticket", r.db.WithContext(ctx).Model(&models.TipoTicket{}).Where("id = ?", item.ID).Updates(updates).Error)
}

// ==================== niveles_prioridad ====================

func (r *CatalogoRepository) ListNivelesPrioridad(ctx context.Context) ([]domain.NivelPrioridad, error) {
	return listCatalogo(r.db, ctx, "nivel de prioridad", func(row models.NivelPrioridad) domain.NivelPrioridad {
		return domain.NivelPrioridad{ID: row.ID, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) GetNivelPrioridadByID(ctx context.Context, id int) (domain.NivelPrioridad, error) {
	return getCatalogoByID(r.db, ctx, id, "nivel de prioridad", func(row models.NivelPrioridad) domain.NivelPrioridad {
		return domain.NivelPrioridad{ID: row.ID, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) CreateNivelPrioridad(ctx context.Context, item *domain.NivelPrioridad) error {
	row := models.NivelPrioridad{Descripcion: item.Descripcion}
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
	return wrapUpdateError("nivel de prioridad", r.db.WithContext(ctx).Model(&models.NivelPrioridad{}).Where("id = ?", item.ID).Updates(updates).Error)
}

// ==================== tipo_tecnico ====================

func (r *CatalogoRepository) ListTiposTecnico(ctx context.Context) ([]domain.TipoTecnico, error) {
	return listCatalogo(r.db, ctx, "tipo técnico", func(row models.TipoTecnico) domain.TipoTecnico {
		return domain.TipoTecnico{ID: row.ID, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) GetTipoTecnicoByID(ctx context.Context, id int) (domain.TipoTecnico, error) {
	return getCatalogoByID(r.db, ctx, id, "tipo técnico", func(row models.TipoTecnico) domain.TipoTecnico {
		return domain.TipoTecnico{ID: row.ID, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) CreateTipoTecnico(ctx context.Context, item *domain.TipoTecnico) error {
	row := models.TipoTecnico{Descripcion: item.Descripcion}
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
	return wrapUpdateError("tipo técnico", r.db.WithContext(ctx).Model(&models.TipoTecnico{}).Where("id = ?", item.ID).Updates(updates).Error)
}

// ==================== departamentos_soporte ====================

func (r *CatalogoRepository) ListDepartamentosSoporte(ctx context.Context) ([]domain.DepartamentoSoporte, error) {
	return listCatalogo(r.db, ctx, "departamento de soporte", func(row models.DepartamentoSoporte) domain.DepartamentoSoporte {
		return domain.DepartamentoSoporte{ID: row.ID, CodDepartamento: row.CodDepartamento, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) GetDepartamentoSoporteByID(ctx context.Context, id int) (domain.DepartamentoSoporte, error) {
	return getCatalogoByID(r.db, ctx, id, "departamento de soporte", func(row models.DepartamentoSoporte) domain.DepartamentoSoporte {
		return domain.DepartamentoSoporte{ID: row.ID, CodDepartamento: row.CodDepartamento, Descripcion: row.Descripcion}
	})
}

func (r *CatalogoRepository) CreateDepartamentoSoporte(ctx context.Context, item *domain.DepartamentoSoporte) error {
	row := models.DepartamentoSoporte{CodDepartamento: item.CodDepartamento, Descripcion: item.Descripcion}
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
	return wrapUpdateError("departamento de soporte", r.db.WithContext(ctx).Model(&models.DepartamentoSoporte{}).Where("id = ?", item.ID).Updates(updates).Error)
}

// ==================== motivos_pausa ====================

func (r *CatalogoRepository) ListMotivosPausa(ctx context.Context) ([]domain.MotivoPausa, error) {
	return listCatalogo(r.db, ctx, "motivo de pausa", func(row models.MotivoPausa) domain.MotivoPausa {
		return domain.MotivoPausa{ID: row.ID, MotivoPausa: row.MotivoPausa, RequiereAutorizacion: derefBool(row.RequiereAutorizacion)}
	})
}

func (r *CatalogoRepository) GetMotivoPausaByID(ctx context.Context, id int) (domain.MotivoPausa, error) {
	return getCatalogoByID(r.db, ctx, id, "motivo de pausa", func(row models.MotivoPausa) domain.MotivoPausa {
		return domain.MotivoPausa{ID: row.ID, MotivoPausa: row.MotivoPausa, RequiereAutorizacion: derefBool(row.RequiereAutorizacion)}
	})
}

func (r *CatalogoRepository) CreateMotivoPausa(ctx context.Context, item *domain.MotivoPausa) error {
	row := models.MotivoPausa{MotivoPausa: item.MotivoPausa, RequiereAutorizacion: &item.RequiereAutorizacion}
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
	return wrapUpdateError("motivo de pausa", r.db.WithContext(ctx).Model(&models.MotivoPausa{}).Where("id = ?", item.ID).Updates(updates).Error)
}

// ==================== tipos_turno ====================

func (r *CatalogoRepository) ListTiposTurno(ctx context.Context) ([]domain.TipoTurno, error) {
	return listCatalogo(r.db, ctx, "tipo de turno", func(row models.TipoTurno) domain.TipoTurno {
		return domain.TipoTurno{ID: row.ID, Nombre: row.Nombre, Descripcion: row.Descripcion, Estado: derefBool(row.Estado)}
	})
}

func (r *CatalogoRepository) GetTipoTurnoByID(ctx context.Context, id int) (domain.TipoTurno, error) {
	return getCatalogoByID(r.db, ctx, id, "tipo de turno", func(row models.TipoTurno) domain.TipoTurno {
		return domain.TipoTurno{ID: row.ID, Nombre: row.Nombre, Descripcion: row.Descripcion, Estado: derefBool(row.Estado)}
	})
}

func (r *CatalogoRepository) CreateTipoTurno(ctx context.Context, item *domain.TipoTurno) error {
	row := models.TipoTurno{Nombre: item.Nombre, Descripcion: item.Descripcion, Estado: &item.Estado}
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
	return wrapUpdateError("tipo de turno", r.db.WithContext(ctx).Model(&models.TipoTurno{}).Where("id = ?", item.ID).Updates(updates).Error)
}

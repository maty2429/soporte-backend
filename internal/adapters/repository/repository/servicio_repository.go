package repository

import (
	"context"

	"gorm.io/gorm"

	dbmodels "soporte/internal/adapters/repository/models"
	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type ServicioRepository struct {
	db *gorm.DB
}

func NewServicioRepository(db *gorm.DB) ports.ServicioRepository {
	if db == nil {
		return unavailableServicioRepository{}
	}
	return &ServicioRepository{db: db}
}

func (r *ServicioRepository) List(ctx context.Context, f ports.ListServiciosFilters) ([]domain.Servicio, int64, error) {
	q := r.db.WithContext(ctx).Model(&dbmodels.Servicio{})

	if f.Edificio != "" {
		q = q.Where("edificio = ?", f.Edificio)
	}
	if f.Piso != nil {
		q = q.Where("piso = ?", *f.Piso)
	}
	if f.Servicios != "" {
		q = q.Where("servicios = ?", f.Servicios)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("(unidades ILIKE ? OR ubicacion ILIKE ?)", like, like)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, wrapListError("servicio", err)
	}

	var rows []dbmodels.Servicio
	if err := q.Order("id ASC").Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, wrapListError("servicio", err)
	}

	items := make([]domain.Servicio, 0, len(rows))
	for _, row := range rows {
		items = append(items, toServicioDomain(row))
	}
	return items, total, nil
}

func (r *ServicioRepository) GetByID(ctx context.Context, id int) (domain.Servicio, error) {
	var row dbmodels.Servicio
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return domain.Servicio{}, wrapDBError("servicio", err)
	}
	return toServicioDomain(row), nil
}

func (r *ServicioRepository) Create(ctx context.Context, item *domain.Servicio) error {
	row := dbmodels.Servicio{
		Edificio:                item.Edificio,
		Piso:                    item.Piso,
		Servicios:               item.Servicios,
		Ubicacion:               item.Ubicacion,
		Unidades:                item.Unidades,
		IDNivelPrioridadDefault: item.IDNivelPrioridadDefault,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("servicio", err)
	}
	item.ID = row.ID
	return nil
}

func (r *ServicioRepository) Update(ctx context.Context, item *domain.Servicio) error {
	row := dbmodels.Servicio{
		ID:                      item.ID,
		Edificio:                item.Edificio,
		Piso:                    item.Piso,
		Servicios:               item.Servicios,
		Ubicacion:               item.Ubicacion,
		Unidades:                item.Unidades,
		IDNivelPrioridadDefault: item.IDNivelPrioridadDefault,
	}
	err := r.db.WithContext(ctx).
		Model(&row).
		Select("Edificio", "Piso", "Servicios", "Ubicacion", "Unidades", "IDNivelPrioridadDefault").
		Updates(&row).Error
	return wrapUpdateError("servicio", err)
}

func toServicioDomain(row dbmodels.Servicio) domain.Servicio {
	return domain.Servicio{
		ID:                      row.ID,
		Edificio:                row.Edificio,
		Piso:                    row.Piso,
		Servicios:               row.Servicios,
		Ubicacion:               row.Ubicacion,
		Unidades:                row.Unidades,
		IDNivelPrioridadDefault: row.IDNivelPrioridadDefault,
	}
}

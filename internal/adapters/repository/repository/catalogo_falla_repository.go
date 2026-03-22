package repository

import (
	"context"

	"gorm.io/gorm"

	dbmodels "soporte/internal/adapters/repository/models"
	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type CatalogoFallaRepository struct {
	db *gorm.DB
}

func NewCatalogoFallaRepository(db *gorm.DB) ports.CatalogoFallaRepository {
	if db == nil {
		return unavailableCatalogoFallaRepository{}
	}
	return &CatalogoFallaRepository{db: db}
}

func (r *CatalogoFallaRepository) List(ctx context.Context) ([]domain.CatalogoFalla, error) {
	var rows []dbmodels.CatalogoFalla
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, wrapListError("catalogo falla", err)
	}
	items := make([]domain.CatalogoFalla, 0, len(rows))
	for _, row := range rows {
		items = append(items, toCatalogoFallaDomain(row))
	}
	return items, nil
}

func (r *CatalogoFallaRepository) GetByID(ctx context.Context, id int) (domain.CatalogoFalla, error) {
	var row dbmodels.CatalogoFalla
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return domain.CatalogoFalla{}, wrapDBError("catalogo falla", err)
	}
	return toCatalogoFallaDomain(row), nil
}

func (r *CatalogoFallaRepository) Create(ctx context.Context, item *domain.CatalogoFalla) error {
	row := dbmodels.CatalogoFalla{
		CodigoFalla:          item.CodigoFalla,
		DescripcionFalla:     item.DescripcionFalla,
		Complejidad:          item.Complejidad,
		RequiereVisitaFisica: &item.RequiereVisitaFisica,
		IDDepartamento:       item.IDDepartamento,
		Categoria:            item.Categoria,
		Subcategoria:         item.Subcategoria,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return wrapCreateError("catalogo falla", err)
	}
	item.ID = row.ID
	return nil
}

func (r *CatalogoFallaRepository) Update(ctx context.Context, item *domain.CatalogoFalla) error {
	updates := map[string]any{}
	if item.CodigoFalla != "" {
		updates["codigo_falla"] = item.CodigoFalla
	}
	if item.DescripcionFalla != "" {
		updates["descripcion_falla"] = item.DescripcionFalla
	}
	if item.Complejidad > 0 {
		updates["complejidad"] = item.Complejidad
	}
	updates["requiere_visita_fisica"] = item.RequiereVisitaFisica
	if item.IDDepartamento != nil {
		updates["id_departamento"] = *item.IDDepartamento
	}
	if item.Categoria != "" {
		updates["categoria"] = item.Categoria
	}
	if item.Subcategoria != "" {
		updates["subcategoria"] = item.Subcategoria
	}
	err := r.db.WithContext(ctx).Model(&dbmodels.CatalogoFalla{}).Where("id = ?", item.ID).Updates(updates).Error
	return wrapUpdateError("catalogo falla", err)
}

func toCatalogoFallaDomain(row dbmodels.CatalogoFalla) domain.CatalogoFalla {
	return domain.CatalogoFalla{
		ID:                       row.ID,
		CodigoFalla:              row.CodigoFalla,
		DescripcionFalla:         row.DescripcionFalla,
		Complejidad:              row.Complejidad,
		TiempoResolucionEstimado: row.TiempoResolucionEstimado,
		RequiereVisitaFisica:     derefBool(row.RequiereVisitaFisica),
		IDDepartamento:           row.IDDepartamento,
		Categoria:                row.Categoria,
		Subcategoria:             row.Subcategoria,
	}
}

package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"

	dbmodels "soporte/internal/adapters/repository/models"
	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

type SolicitanteRepository struct {
	db *gorm.DB
}

func NewSolicitanteRepository(db *gorm.DB) ports.SolicitanteRepository {
	if db == nil {
		return unavailableSolicitanteRepository{}
	}

	return &SolicitanteRepository{db: db}
}

func (r *SolicitanteRepository) List(ctx context.Context, filters ports.ListSolicitantesFilters) ([]domain.Solicitante, int64, error) {
	var (
		rows  []dbmodels.Solicitante
		total int64
	)

	query := r.db.WithContext(ctx).Model(&dbmodels.Solicitante{})

	if filters.Search != "" {
		like := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(nombre_completo) LIKE ? OR rut LIKE ?", like, like)
	}

	if filters.Estado != nil {
		query = query.Where("estado = ?", *filters.Estado)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapListError("solicitantes", err)
	}

	if err := query.
		Order("id ASC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&rows).Error; err != nil {
		return nil, 0, wrapListError("solicitantes", err)
	}

	items := make([]domain.Solicitante, 0, len(rows))
	for _, row := range rows {
		items = append(items, toSolicitanteDomain(row))
	}

	return items, total, nil
}

func (r *SolicitanteRepository) GetByID(ctx context.Context, id int) (domain.Solicitante, error) {
	var row dbmodels.Solicitante

	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return domain.Solicitante{}, wrapDBError("solicitante", err)
	}

	return toSolicitanteDomain(row), nil
}

func (r *SolicitanteRepository) GetByRutDV(ctx context.Context, rut, dv string) (domain.Solicitante, error) {
	var row dbmodels.Solicitante

	if err := r.db.WithContext(ctx).
		Preload("Servicio").
		Where("rut = ? AND dv = ?", rut, dv).
		Take(&row).Error; err != nil {
		return domain.Solicitante{}, wrapDBError("solicitante", err)
	}

	return toSolicitanteDomain(row), nil
}

func (r *SolicitanteRepository) Create(ctx context.Context, sol *domain.Solicitante) error {
	row := toSolicitanteToDB(*sol)

	if err := r.db.WithContext(ctx).
		Select("IDServicio", "Correo", "Rut", "Dv", "NombreCompleto", "Anexo", "Estado").
		Create(&row).Error; err != nil {
		return wrapCreateError("solicitante", err)
	}

	*sol = toSolicitanteDomain(row)
	return nil
}

func (r *SolicitanteRepository) Update(ctx context.Context, sol *domain.Solicitante) error {
	row := toSolicitanteToDB(*sol)

	err := r.db.WithContext(ctx).
		Model(&row).
		Select("IDServicio", "Correo", "Rut", "Dv", "NombreCompleto", "Anexo", "Estado").
		Updates(&row).Error

	return wrapUpdateError("solicitante", err)
}

func toSolicitanteDomain(row dbmodels.Solicitante) domain.Solicitante {
	var correo string
	if row.Correo != nil {
		correo = *row.Correo
	}
	var estado bool
	if row.Estado != nil {
		estado = *row.Estado
	}
	return domain.Solicitante{
		ID:             row.ID,
		IDServicio:     row.IDServicio,
		Servicio:       toServicioPtr(row.Servicio),
		Correo:         correo,
		Rut:            row.Rut,
		Dv:             row.Dv,
		NombreCompleto: row.NombreCompleto,
		Anexo:          row.Anexo,
		Estado:         estado,
	}
}

func toServicioPtr(row *dbmodels.Servicio) *domain.Servicio {
	if row == nil {
		return nil
	}

	return &domain.Servicio{
		ID:                      row.ID,
		Edificio:                row.Edificio,
		Piso:                    row.Piso,
		Servicios:               row.Servicios,
		Ubicacion:               row.Ubicacion,
		Unidades:                row.Unidades,
		IDNivelPrioridadDefault: row.IDNivelPrioridadDefault,
	}
}

func toSolicitanteToDB(sol domain.Solicitante) dbmodels.Solicitante {
	var correo *string
	if sol.Correo != "" {
		correo = &sol.Correo
	}
	return dbmodels.Solicitante{
		ID:             sol.ID,
		IDServicio:     sol.IDServicio,
		Correo:         correo,
		Rut:            sol.Rut,
		Dv:             sol.Dv,
		NombreCompleto: sol.NombreCompleto,
		Anexo:          sol.Anexo,
		Estado:         &sol.Estado,
	}
}

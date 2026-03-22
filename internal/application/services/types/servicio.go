package types

import "soporte/internal/core/domain"

type ListServiciosQuery struct {
	Edificio  string
	Piso      *int
	Servicios string
	Search    string
	Limit     int
	Offset    int
}

type ListServiciosResult struct {
	Items  []domain.Servicio
	Total  int64
	Limit  int
	Offset int
}

type CreateServicioCommand struct {
	Edificio                string
	Piso                    int
	Servicios               string
	Ubicacion               string
	Unidades                string
	IDNivelPrioridadDefault *int
}

type UpdateServicioCommand struct {
	ID                      int
	Edificio                *string
	Piso                    *int
	Servicios               *string
	Ubicacion               *string
	Unidades                *string
	IDNivelPrioridadDefault *int
}

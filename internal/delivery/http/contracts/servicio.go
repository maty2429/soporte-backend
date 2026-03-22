package contracts

import "soporte/internal/core/domain"

type ListServiciosQuery struct {
	Edificio  string `form:"edificio"  binding:"omitempty"`
	Piso      *int   `form:"piso"      binding:"omitempty"`
	Servicios string `form:"servicios" binding:"omitempty"`
	Search    string `form:"search"    binding:"omitempty"`
	Limit     int    `form:"limit"     binding:"omitempty,gt=0,lte=100"`
	Offset    int    `form:"offset"    binding:"omitempty,gte=0"`
}

type CreateServicioRequest struct {
	Edificio                string `json:"edificio"`
	Piso                    int    `json:"piso"`
	Servicios               string `json:"servicios"`
	Ubicacion               string `json:"ubicacion"`
	Unidades                string `json:"unidades"`
	IDNivelPrioridadDefault *int   `json:"id_nivel_prioridad_default" binding:"omitempty,gt=0"`
}

type UpdateServicioRequest struct {
	Edificio                *string `json:"edificio"                   binding:"omitempty"`
	Piso                    *int    `json:"piso"                       binding:"omitempty"`
	Servicios               *string `json:"servicios"                  binding:"omitempty"`
	Ubicacion               *string `json:"ubicacion"                  binding:"omitempty"`
	Unidades                *string `json:"unidades"                   binding:"omitempty"`
	IDNivelPrioridadDefault *int    `json:"id_nivel_prioridad_default" binding:"omitempty,gt=0"`
}

type ServicioResponse struct {
	ID                      int    `json:"id"`
	Edificio                string `json:"edificio"`
	Piso                    int    `json:"piso"`
	Servicios               string `json:"servicios"`
	Ubicacion               string `json:"ubicacion"`
	Unidades                string `json:"unidades"`
	IDNivelPrioridadDefault *int   `json:"id_nivel_prioridad_default"`
}

func NewServicioResponse(s domain.Servicio) ServicioResponse {
	return ServicioResponse{
		ID:                      s.ID,
		Edificio:                s.Edificio,
		Piso:                    s.Piso,
		Servicios:               s.Servicios,
		Ubicacion:               s.Ubicacion,
		Unidades:                s.Unidades,
		IDNivelPrioridadDefault: s.IDNivelPrioridadDefault,
	}
}

func NewServiciosResponse(items []domain.Servicio) []ServicioResponse {
	out := make([]ServicioResponse, 0, len(items))
	for _, s := range items {
		out = append(out, NewServicioResponse(s))
	}
	return out
}

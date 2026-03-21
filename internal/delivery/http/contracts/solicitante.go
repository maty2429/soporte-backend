package contracts

import "soporte/internal/core/domain"

type ServicioResponse struct {
	ID                      int    `json:"id"`
	Edificio                string `json:"edificio"`
	Piso                    int    `json:"piso"`
	Servicios               string `json:"servicios"`
	Ubicacion               string `json:"ubicacion"`
	IDNivelPrioridadDefault *int   `json:"id_nivel_prioridad_default,omitempty"`
}

type ListSolicitantesQuery struct {
	Limit  int    `form:"limit"  binding:"omitempty,gt=0,lte=100"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
	Search string `form:"q"      binding:"omitempty,max=100"`
	Estado *bool  `form:"estado" binding:"omitempty"`
}

type CreateSolicitanteRequest struct {
	IDServicio     *int   `json:"id_servicio"     binding:"omitempty,gt=0"`
	Correo         string `json:"correo"          binding:"omitempty,email,max=100"`
	Rut            string `json:"rut"             binding:"required,max=10"`
	Dv             string `json:"dv"              binding:"required,len=1"`
	NombreCompleto string `json:"nombre_completo" binding:"required"`
	Anexo          *int   `json:"anexo"           binding:"omitempty,gt=0"`
	Estado         *bool  `json:"estado"          binding:"omitempty"`
}

type UpdateSolicitanteRequest struct {
	IDServicio     *int    `json:"id_servicio"     binding:"omitempty,gt=0"`
	Correo         *string `json:"correo"          binding:"omitempty,email,max=100"`
	Rut            *string `json:"rut"             binding:"omitempty,max=10"`
	Dv             *string `json:"dv"              binding:"omitempty,len=1"`
	NombreCompleto *string `json:"nombre_completo" binding:"omitempty"`
	Anexo          *int    `json:"anexo"           binding:"omitempty,gt=0"`
	Estado         *bool   `json:"estado"          binding:"omitempty"`
}

type SolicitanteResponse struct {
	ID             int               `json:"id"`
	IDServicio     *int              `json:"id_servicio,omitempty"`
	Servicio       *ServicioResponse `json:"servicio,omitempty"`
	Correo         string            `json:"correo"`
	Rut            string            `json:"rut"`
	Dv             string            `json:"dv"`
	NombreCompleto string            `json:"nombre_completo"`
	Anexo          *int              `json:"anexo,omitempty"`
	Estado         bool              `json:"estado"`
}

func NewSolicitanteResponse(sol domain.Solicitante) SolicitanteResponse {
	return SolicitanteResponse{
		ID:             sol.ID,
		IDServicio:     sol.IDServicio,
		Servicio:       newServicioResponse(sol.Servicio),
		Correo:         sol.Correo,
		Rut:            sol.Rut,
		Dv:             sol.Dv,
		NombreCompleto: sol.NombreCompleto,
		Anexo:          sol.Anexo,
		Estado:         sol.Estado,
	}
}

func newServicioResponse(servicio *domain.Servicio) *ServicioResponse {
	if servicio == nil {
		return nil
	}

	return &ServicioResponse{
		ID:                      servicio.ID,
		Edificio:                servicio.Edificio,
		Piso:                    servicio.Piso,
		Servicios:               servicio.Servicios,
		Ubicacion:               servicio.Ubicacion,
		IDNivelPrioridadDefault: servicio.IDNivelPrioridadDefault,
	}
}

type SolicitanteCreatedResponse struct {
	ID int `json:"id"`
}

func NewSolicitanteCreatedResponse(id int) SolicitanteCreatedResponse {
	return SolicitanteCreatedResponse{ID: id}
}

func NewSolicitantesResponse(items []domain.Solicitante) []SolicitanteResponse {
	response := make([]SolicitanteResponse, 0, len(items))
	for _, sol := range items {
		response = append(response, NewSolicitanteResponse(sol))
	}
	return response
}

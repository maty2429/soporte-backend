package contracts

import "soporte/internal/core/domain"

type CreateCatalogoFallaRequest struct {
	CodigoFalla          string `json:"codigo_falla"           binding:"required"`
	DescripcionFalla     string `json:"descripcion_falla"      binding:"required"`
	Complejidad          int    `json:"complejidad"            binding:"required,gte=1,lte=10"`
	RequiereVisitaFisica bool   `json:"requiere_visita_fisica"`
	IDDepartamento       *int   `json:"id_departamento"        binding:"omitempty,gt=0"`
	Categoria            string `json:"categoria"              binding:"omitempty"`
	Subcategoria         string `json:"subcategoria"           binding:"omitempty"`
}

type UpdateCatalogoFallaRequest struct {
	CodigoFalla          *string `json:"codigo_falla"           binding:"omitempty"`
	DescripcionFalla     *string `json:"descripcion_falla"      binding:"omitempty"`
	Complejidad          *int    `json:"complejidad"            binding:"omitempty,gte=1,lte=10"`
	RequiereVisitaFisica *bool   `json:"requiere_visita_fisica" binding:"omitempty"`
	IDDepartamento       *int    `json:"id_departamento"        binding:"omitempty,gt=0"`
	Categoria            *string `json:"categoria"              binding:"omitempty"`
	Subcategoria         *string `json:"subcategoria"           binding:"omitempty"`
}

type CatalogoFallaResponse struct {
	ID                       int    `json:"id"`
	CodigoFalla              string `json:"codigo_falla"`
	DescripcionFalla         string `json:"descripcion_falla"`
	Complejidad              int    `json:"complejidad"`
	TiempoResolucionEstimado string `json:"tiempo_resolucion_estimado,omitempty"`
	RequiereVisitaFisica     bool   `json:"requiere_visita_fisica"`
	IDDepartamento           *int   `json:"id_departamento"`
	Categoria                string `json:"categoria,omitempty"`
	Subcategoria             string `json:"subcategoria,omitempty"`
}

func NewCatalogoFallaResponse(f domain.CatalogoFalla) CatalogoFallaResponse {
	return CatalogoFallaResponse{
		ID:                       f.ID,
		CodigoFalla:              f.CodigoFalla,
		DescripcionFalla:         f.DescripcionFalla,
		Complejidad:              f.Complejidad,
		TiempoResolucionEstimado: f.TiempoResolucionEstimado,
		RequiereVisitaFisica:     f.RequiereVisitaFisica,
		IDDepartamento:           f.IDDepartamento,
		Categoria:                f.Categoria,
		Subcategoria:             f.Subcategoria,
	}
}

func NewCatalogoFallasResponse(items []domain.CatalogoFalla) []CatalogoFallaResponse {
	out := make([]CatalogoFallaResponse, 0, len(items))
	for _, f := range items {
		out = append(out, NewCatalogoFallaResponse(f))
	}
	return out
}

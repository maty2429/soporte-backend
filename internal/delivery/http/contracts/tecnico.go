package contracts

import "soporte/internal/core/domain"

type ListTecnicosQuery struct {
	Limit                 int    `form:"limit"                  binding:"omitempty,gt=0,lte=100"`
	Offset                int    `form:"offset"                 binding:"omitempty,gte=0"`
	Search                string `form:"q"                      binding:"omitempty,max=100"`
	Estado                *bool  `form:"estado"                 binding:"omitempty"`
	IDTipoTecnico         int    `form:"id_tipo_tecnico"        binding:"omitempty,gt=0"`
	IDDepartamentoSoporte int    `form:"id_departamento_soporte" binding:"omitempty,gt=0"`
}

type CreateTecnicoRequest struct {
	Rut                   string `json:"rut"                     binding:"required,max=10"`
	Dv                    string `json:"dv"                      binding:"required,len=1"`
	NombreCompleto        string `json:"nombre_completo"         binding:"required"`
	IDTipoTecnico         *int   `json:"id_tipo_tecnico"         binding:"omitempty,gt=0"`
	IDDepartamentoSoporte *int   `json:"id_departamento_soporte" binding:"omitempty,gt=0"`
	IDTipoTurno           *int   `json:"id_tipo_turno"           binding:"omitempty,gt=0"`
	Estado                *bool  `json:"estado"                  binding:"omitempty"`
}

type UpdateTecnicoRequest struct {
	Rut                   *string `json:"rut"                     binding:"omitempty,max=10"`
	Dv                    *string `json:"dv"                      binding:"omitempty,len=1"`
	NombreCompleto        *string `json:"nombre_completo"         binding:"omitempty"`
	IDTipoTecnico         *int    `json:"id_tipo_tecnico"         binding:"omitempty,gt=0"`
	IDDepartamentoSoporte *int    `json:"id_departamento_soporte" binding:"omitempty,gt=0"`
	IDTipoTurno           *int    `json:"id_tipo_turno"           binding:"omitempty,gt=0"`
	Estado                *bool   `json:"estado"                  binding:"omitempty"`
}

type TecnicoResponse struct {
	ID                    int                                `json:"id"`
	Rut                   string                             `json:"rut"`
	Dv                    string                             `json:"dv"`
	NombreCompleto        string                             `json:"nombre_completo"`
	IDTipoTecnico         *int                               `json:"id_tipo_tecnico,omitempty"`
	IDDepartamentoSoporte *int                               `json:"id_departamento_soporte,omitempty"`
	IDTipoTurno           *int                               `json:"id_tipo_turno,omitempty"`
	Estado                bool                               `json:"estado"`
	CreatedAt             string                             `json:"created_at"`
	UpdatedAt             string                             `json:"updated_at"`
	DepartamentoSoporte   *TicketDepartamentoSoporteResponse `json:"departamento_soporte,omitempty"`
}

func NewTecnicoResponse(t domain.Tecnico) TecnicoResponse {
	const timeFmt = "2006-01-02T15:04:05Z07:00"
	return TecnicoResponse{
		ID:                    t.ID,
		Rut:                   t.Rut,
		Dv:                    t.Dv,
		NombreCompleto:        t.NombreCompleto,
		IDTipoTecnico:         t.IDTipoTecnico,
		IDDepartamentoSoporte: t.IDDepartamentoSoporte,
		IDTipoTurno:           t.IDTipoTurno,
		Estado:                t.Estado,
		CreatedAt:             t.CreatedAt.Format(timeFmt),
		UpdatedAt:             t.UpdatedAt.Format(timeFmt),
		DepartamentoSoporte:   newTicketDepartamentoSoporteResponse(t.DepartamentoSoporte),
	}
}

func NewTecnicosResponse(items []domain.Tecnico) []TecnicoResponse {
	out := make([]TecnicoResponse, 0, len(items))
	for _, t := range items {
		out = append(out, NewTecnicoResponse(t))
	}
	return out
}

type TecnicoCreatedResponse struct {
	ID int `json:"id"`
}

func NewTecnicoCreatedResponse(id int) TecnicoCreatedResponse {
	return TecnicoCreatedResponse{ID: id}
}

// --- Configuración Horarios Turno ---

type CreateHorarioTurnoRequest struct {
	IDTipoTurno int    `json:"id_tipo_turno" binding:"required,gt=0"`
	DiaSemana   int    `json:"dia_semana"    binding:"gte=0,lte=6"`
	HoraInicio  string `json:"hora_inicio"   binding:"required"`
	HoraFin     string `json:"hora_fin"      binding:"required"`
}

type UpdateHorarioTurnoRequest struct {
	IDTipoTurno *int    `json:"id_tipo_turno" binding:"omitempty,gt=0"`
	DiaSemana   *int    `json:"dia_semana"    binding:"omitempty,gte=0,lte=6"`
	HoraInicio  *string `json:"hora_inicio"   binding:"omitempty"`
	HoraFin     *string `json:"hora_fin"      binding:"omitempty"`
}

type HorarioTurnoResponse struct {
	ID          int    `json:"id"`
	IDTipoTurno int    `json:"id_tipo_turno"`
	DiaSemana   int    `json:"dia_semana"`
	HoraInicio  string `json:"hora_inicio"`
	HoraFin     string `json:"hora_fin"`
}

func NewHorarioTurnoResponse(h domain.ConfiguracionHorarioTurno) HorarioTurnoResponse {
	return HorarioTurnoResponse{
		ID:          h.ID,
		IDTipoTurno: h.IDTipoTurno,
		DiaSemana:   h.DiaSemana,
		HoraInicio:  h.HoraInicio,
		HoraFin:     h.HoraFin,
	}
}

func NewHorariosTurnoResponse(items []domain.ConfiguracionHorarioTurno) []HorarioTurnoResponse {
	out := make([]HorarioTurnoResponse, 0, len(items))
	for _, h := range items {
		out = append(out, NewHorarioTurnoResponse(h))
	}
	return out
}

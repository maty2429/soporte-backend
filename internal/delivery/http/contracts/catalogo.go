package contracts

// --- tipos_ticket ---

type CreateTipoTicketRequest struct {
	CodTipoTicket string `json:"cod_tipo_ticket"`
	Descripcion   string `json:"descripcion"`
}

type UpdateTipoTicketRequest struct {
	CodTipoTicket *string `json:"cod_tipo_ticket"`
	Descripcion   *string `json:"descripcion"`
}

// --- niveles_prioridad, tipo_tecnico (solo descripcion) ---

type CreateDescripcionRequest struct {
	Descripcion string `json:"descripcion"`
}

type UpdateDescripcionRequest struct {
	Descripcion *string `json:"descripcion"`
}

// --- departamentos_soporte ---

type CreateDepartamentoSoporteRequest struct {
	CodDepartamento string `json:"cod_departamento"`
	Descripcion     string `json:"descripcion"`
}

type UpdateDepartamentoSoporteRequest struct {
	CodDepartamento *string `json:"cod_departamento"`
	Descripcion     *string `json:"descripcion"`
}

// --- motivos_pausa ---

type CreateMotivoPausaRequest struct {
	MotivoPausa          string `json:"motivo_pausa"`
	RequiereAutorizacion *bool  `json:"requiere_autorizacion"`
}

type UpdateMotivoPausaRequest struct {
	MotivoPausa          *string `json:"motivo_pausa"`
	RequiereAutorizacion *bool   `json:"requiere_autorizacion"`
}

// --- tipos_turno ---

type CreateTipoTurnoRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion *string `json:"descripcion"`
	Estado      *bool   `json:"estado"`
}

type UpdateTipoTurnoRequest struct {
	Nombre      *string `json:"nombre"`
	Descripcion *string `json:"descripcion"`
	Estado      *bool   `json:"estado"`
}

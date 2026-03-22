package contracts

import "soporte/internal/core/domain"

type CreateTicketRequest struct {
	IDSolicitante         int     `json:"id_solicitante"          binding:"required,gt=0"`
	IDServicio            int     `json:"id_servicio"             binding:"required,gt=0"`
	IDTipoTicket          int     `json:"id_tipo_ticket"          binding:"required,gt=0"`
	IDNivelPrioridad      int     `json:"id_nivel_prioridad"      binding:"required,gt=0"`
	IDDepartamentoSoporte int     `json:"id_departamento_soporte" binding:"required,gt=0"`
	Critico               bool    `json:"critico"`
	DetalleFallaReportada string  `json:"detalle_falla_reportada" binding:"required"`
	UbicacionObs          *string `json:"ubicacion_obs,omitempty"`
}

type TicketCreatedResponse struct {
	NroTicket string `json:"nro_ticket"`
}

func NewTicketCreatedResponse(ticket domain.Ticket) TicketCreatedResponse {
	return TicketCreatedResponse{
		NroTicket: ticket.NroTicket,
	}
}

type AssignTicketRequest struct {
	IDTecnicoAsignado int `json:"id_tecnico_asignado" binding:"required,gt=0"`
	IDCatalogoFalla   int `json:"id_catalogo_falla"   binding:"required,gt=0"`
	IDNivelPrioridad  int `json:"id_nivel_prioridad"  binding:"required,gt=0"`
}

type AssignTicketResponse struct {
	ID                    int    `json:"id"`
	NroTicket             string `json:"nro_ticket"`
	CodEstadoTicket       string `json:"cod_estado_ticket"`
	IDTecnicoAsignado     *int   `json:"id_tecnico_asignado"`
	IDCatalogoFalla       *int   `json:"id_catalogo_falla"`
	IDNivelPrioridad      *int   `json:"id_nivel_prioridad"`
}

func NewAssignTicketResponse(ticket domain.Ticket) AssignTicketResponse {
	return AssignTicketResponse{
		ID:                ticket.ID,
		NroTicket:         ticket.NroTicket,
		CodEstadoTicket:   ticket.CodEstadoTicket,
		IDTecnicoAsignado: ticket.IDTecnicoAsignado,
		IDCatalogoFalla:   ticket.IDCatalogoFalla,
		IDNivelPrioridad:  ticket.IDNivelPrioridad,
	}
}

type UpdateTicketRequest struct {
	DetalleFallaReportada *string `json:"detalle_falla_reportada,omitempty"`
	UbicacionObs          *string `json:"ubicacion_obs,omitempty"`
	Critico               *bool   `json:"critico,omitempty"`
	IDTipoTicket          *int    `json:"id_tipo_ticket,omitempty"`
	IDDepartamentoSoporte *int    `json:"id_departamento_soporte,omitempty"`
	IDServicio            *int    `json:"id_servicio,omitempty"`
}

type ListTicketsQuery struct {
	Limit                 int    `form:"limit"                   binding:"omitempty,gt=0,lte=100"`
	Offset                int    `form:"offset"                  binding:"omitempty,gte=0"`
	CodEstadoTicket       string `form:"estado"                  binding:"omitempty"`
	IDTecnicoAsignado     int    `form:"id_tecnico"              binding:"omitempty,gt=0"`
	IDSolicitante         int    `form:"id_solicitante"          binding:"omitempty,gt=0"`
	IDDepartamentoSoporte int    `form:"id_departamento"         binding:"omitempty,gt=0"`
	Critico               *bool  `form:"critico"                 binding:"omitempty"`
}

type TicketResponse struct {
	ID                    int     `json:"id"`
	NroTicket             string  `json:"nro_ticket"`
	IDSolicitante         int     `json:"id_solicitante"`
	IDTecnicoAsignado     *int    `json:"id_tecnico_asignado"`
	IDServicio            *int    `json:"id_servicio"`
	IDTipoTicket          int     `json:"id_tipo_ticket"`
	CodEstadoTicket       string  `json:"cod_estado_ticket"`
	IDNivelPrioridad      *int    `json:"id_nivel_prioridad"`
	IDCatalogoFalla       *int    `json:"id_catalogo_falla"`
	IDDepartamentoSoporte *int    `json:"id_departamento_soporte"`
	Critico               bool    `json:"critico"`
	DetalleFallaReportada string  `json:"detalle_falla_reportada"`
	UbicacionObs          string  `json:"ubicacion_obs"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	FechaInicioTrabajo    *string `json:"fecha_inicio_trabajo"`
	FechaFinTrabajo       *string `json:"fecha_fin_trabajo"`
}

func NewTicketResponse(t domain.Ticket) TicketResponse {
	const timeFmt = "2006-01-02T15:04:05Z07:00"
	resp := TicketResponse{
		ID:                    t.ID,
		NroTicket:             t.NroTicket,
		IDSolicitante:         t.IDSolicitante,
		IDTecnicoAsignado:     t.IDTecnicoAsignado,
		IDServicio:            t.IDServicio,
		IDTipoTicket:          t.IDTipoTicket,
		CodEstadoTicket:       t.CodEstadoTicket,
		IDNivelPrioridad:      t.IDNivelPrioridad,
		IDCatalogoFalla:       t.IDCatalogoFalla,
		IDDepartamentoSoporte: t.IDDepartamentoSoporte,
		Critico:               t.Critico,
		DetalleFallaReportada: t.DetalleFallaReportada,
		UbicacionObs:          t.UbicacionObs,
		CreatedAt:             t.CreatedAt.Format(timeFmt),
		UpdatedAt:             t.UpdatedAt.Format(timeFmt),
	}
	if t.FechaInicioTrabajo != nil {
		s := t.FechaInicioTrabajo.Format(timeFmt)
		resp.FechaInicioTrabajo = &s
	}
	if t.FechaFinTrabajo != nil {
		s := t.FechaFinTrabajo.Format(timeFmt)
		resp.FechaFinTrabajo = &s
	}
	return resp
}

func NewTicketsResponse(items []domain.Ticket) []TicketResponse {
	out := make([]TicketResponse, 0, len(items))
	for _, t := range items {
		out = append(out, NewTicketResponse(t))
	}
	return out
}

func NewBitacorasResponse(items []domain.BitacoraTicket) []BitacoraResponse {
	out := make([]BitacoraResponse, 0, len(items))
	for _, b := range items {
		out = append(out, NewBitacoraResponse(b))
	}
	return out
}

type CreatePausaRequest struct {
	IDTecnicoPausa int `json:"id_tecnico_pausa" binding:"required,gt=0"`
	IDMotivoPausa  int `json:"id_motivo_pausa"  binding:"required,gt=0"`
}

type PausaResponse struct {
	ID               int    `json:"id"`
	IDTicket         int    `json:"id_ticket"`
	IDTecnicoPausa   int    `json:"id_tecnico_pausa"`
	EstadoPausa      string `json:"estado_pausa"`
	IDMotivoPausa    int    `json:"id_motivo_pausa"`
	FechaInicioPausa string `json:"fecha_inicio_pausa"`
}

func NewPausaResponse(p domain.TicketPausa) PausaResponse {
	return PausaResponse{
		ID:               p.ID,
		IDTicket:         p.IDTicket,
		IDTecnicoPausa:   p.IDTecnicoPausa,
		EstadoPausa:      p.EstadoPausa,
		IDMotivoPausa:    p.IDMotivoPausa,
		FechaInicioPausa: p.FechaInicioPausa.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type ResolverPausaRequest struct {
	EstadoPausa         string `json:"estado_pausa"          binding:"required"`
	IDTecnicoAutorizado int    `json:"id_tecnico_autorizado" binding:"required,gt=0"`
}

type ListPausasQuery struct {
	Limit  int    `form:"limit"  binding:"omitempty,gt=0,lte=100"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
	Estado string `form:"estado" binding:"omitempty"`
}

type PausaDetalleResponse struct {
	ID                  int     `json:"id"`
	IDTicket            int     `json:"id_ticket"`
	IDTecnicoPausa      int     `json:"id_tecnico_pausa"`
	EstadoPausa         string  `json:"estado_pausa"`
	IDMotivoPausa       int     `json:"id_motivo_pausa"`
	FechaInicioPausa    string  `json:"fecha_inicio_pausa"`
	FechaFinPausa       *string `json:"fecha_fin_pausa"`
	IDTecnicoAutorizado *int    `json:"id_tecnico_autorizado"`
	FechaResolucion     *string `json:"fecha_resolucion"`
}

func NewPausaDetalleResponse(p domain.TicketPausa) PausaDetalleResponse {
	const timeFmt = "2006-01-02T15:04:05Z07:00"
	resp := PausaDetalleResponse{
		ID:                  p.ID,
		IDTicket:            p.IDTicket,
		IDTecnicoPausa:      p.IDTecnicoPausa,
		EstadoPausa:         p.EstadoPausa,
		IDMotivoPausa:       p.IDMotivoPausa,
		FechaInicioPausa:    p.FechaInicioPausa.Format(timeFmt),
		IDTecnicoAutorizado: p.IDTecnicoAutorizado,
	}
	if p.FechaFinPausa != nil {
		s := p.FechaFinPausa.Format(timeFmt)
		resp.FechaFinPausa = &s
	}
	if p.FechaResolucion != nil {
		s := p.FechaResolucion.Format(timeFmt)
		resp.FechaResolucion = &s
	}
	return resp
}

func NewPausasDetalleResponse(items []domain.TicketPausa) []PausaDetalleResponse {
	out := make([]PausaDetalleResponse, 0, len(items))
	for _, p := range items {
		out = append(out, NewPausaDetalleResponse(p))
	}
	return out
}

type ReanudarTicketRequest struct {
	IDTecnicoPausa int `json:"id_tecnico_pausa" binding:"required,gt=0"`
}

type CloseTicketRequest struct {
	IDSolicitante int    `json:"id_solicitante" binding:"required,gt=0"`
	Nota          int    `json:"nota"           binding:"required,min=1,max=5"`
	Comentarios   string `json:"comentarios"`
	Observacion   string `json:"observacion"    binding:"required"`
}

type ChangeEstadoRequest struct {
	CodEstadoTicket string `json:"cod_estado_ticket" binding:"required"`
	RutResponsable  string `json:"rut_responsable"   binding:"required"`
}

type CreateTraspasoRequest struct {
	IDTecnicoOrigen  int    `json:"id_tecnico_origen"  binding:"required,gt=0"`
	IDTecnicoDestino int    `json:"id_tecnico_destino" binding:"required,gt=0"`
	Motivo           string `json:"motivo"             binding:"required"`
}

type ResolverTraspasoRequest struct {
	EstadoTraspaso       string `json:"estado_traspaso"       binding:"required"`
	ComentarioResolucion string `json:"comentario_resolucion"`
}

type ListTraspasosQuery struct {
	Limit  int    `form:"limit"  binding:"omitempty,gt=0,lte=100"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
	Estado string `form:"estado" binding:"omitempty"`
}

type TraspasoResponse struct {
	ID                   int     `json:"id"`
	IDTicket             int     `json:"id_ticket"`
	IDTecnicoOrigen      int     `json:"id_tecnico_origen"`
	IDTecnicoDestino     int     `json:"id_tecnico_destino"`
	EstadoTraspaso       string  `json:"estado_traspaso"`
	Motivo               string  `json:"motivo"`
	ComentarioResolucion string  `json:"comentario_resolucion,omitempty"`
	FechaSolicitud       string  `json:"fecha_solicitud"`
	FechaResolucion      *string `json:"fecha_resolucion"`
}

func NewTraspasoResponse(t domain.TicketTraspaso) TraspasoResponse {
	const timeFmt = "2006-01-02T15:04:05Z07:00"
	resp := TraspasoResponse{
		ID:                   t.ID,
		IDTicket:             t.IDTicket,
		IDTecnicoOrigen:      t.IDTecnicoOrigen,
		IDTecnicoDestino:     t.IDTecnicoDestino,
		EstadoTraspaso:       t.EstadoTraspaso,
		Motivo:               t.Motivo,
		ComentarioResolucion: t.ComentarioResolucion,
		FechaSolicitud:       t.FechaSolicitud.Format(timeFmt),
	}
	if t.FechaResolucion != nil {
		s := t.FechaResolucion.Format(timeFmt)
		resp.FechaResolucion = &s
	}
	return resp
}

func NewTraspasosResponse(items []domain.TicketTraspaso) []TraspasoResponse {
	out := make([]TraspasoResponse, 0, len(items))
	for _, t := range items {
		out = append(out, NewTraspasoResponse(t))
	}
	return out
}

type CreateBitacoraRequest struct {
	RutAutor   string `json:"rut_autor"   binding:"required"`
	Comentario string `json:"comentario"  binding:"required"`
}

type BitacoraResponse struct {
	ID            int    `json:"id"`
	IDTicket      int    `json:"id_ticket"`
	RutAutor      string `json:"rut_autor"`
	Comentario    string `json:"comentario"`
	FechaRegistro string `json:"fecha_registro"`
}

func NewBitacoraResponse(b domain.BitacoraTicket) BitacoraResponse {
	return BitacoraResponse{
		ID:            b.ID,
		IDTicket:      b.IDTicket,
		RutAutor:      b.RutAutor,
		Comentario:    b.Comentario,
		FechaRegistro: b.FechaRegistro.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type TrazabilidadResponse struct {
	ID                 int    `json:"id"`
	IDTicket           int    `json:"id_ticket"`
	CodEstadoTicket    string `json:"cod_estado_ticket"`
	DescripcionEstado  string `json:"descripcion_estado"`
	RutResponsable     string `json:"rut_responsable"`
	FechaTrazabilidad  string `json:"fecha_trazabilidad"`
}

type TicketDetalleResponse struct {
	ID                    int                     `json:"id"`
	NroTicket             string                  `json:"nro_ticket"`
	IDSolicitante         int                     `json:"id_solicitante"`
	IDTecnicoAsignado     *int                    `json:"id_tecnico_asignado"`
	IDServicio            *int                    `json:"id_servicio"`
	IDTipoTicket          int                     `json:"id_tipo_ticket"`
	CodEstadoTicket       string                  `json:"cod_estado_ticket"`
	IDNivelPrioridad      *int                    `json:"id_nivel_prioridad"`
	IDCatalogoFalla       *int                    `json:"id_catalogo_falla"`
	IDDepartamentoSoporte *int                    `json:"id_departamento_soporte"`
	Critico               bool                    `json:"critico"`
	DetalleFallaReportada string                  `json:"detalle_falla_reportada"`
	UbicacionObs          string                  `json:"ubicacion_obs"`
	CreatedAt             string                  `json:"created_at"`
	UpdatedAt             string                  `json:"updated_at"`
	FechaInicioTrabajo    *string                 `json:"fecha_inicio_trabajo"`
	FechaFinTrabajo       *string                 `json:"fecha_fin_trabajo"`
	Trazabilidad          []TrazabilidadResponse  `json:"trazabilidad"`
	Bitacora              []BitacoraResponse      `json:"bitacora"`
}

func NewTicketDetalleResponse(d domain.TicketDetalle) TicketDetalleResponse {
	const timeFmt = "2006-01-02T15:04:05Z07:00"

	trazabilidad := make([]TrazabilidadResponse, 0, len(d.Trazabilidad))
	for _, t := range d.Trazabilidad {
		trazabilidad = append(trazabilidad, TrazabilidadResponse{
			ID:                t.ID,
			IDTicket:          t.IDTicket,
			CodEstadoTicket:   t.CodEstadoTicket,
			DescripcionEstado: t.DescripcionEstado,
			RutResponsable:    t.RutResponsable,
			FechaTrazabilidad: t.FechaTrazabilidad.Format(timeFmt),
		})
	}

	bitacora := make([]BitacoraResponse, 0, len(d.Bitacora))
	for _, b := range d.Bitacora {
		bitacora = append(bitacora, NewBitacoraResponse(b))
	}

	resp := TicketDetalleResponse{
		ID:                    d.Ticket.ID,
		NroTicket:             d.Ticket.NroTicket,
		IDSolicitante:         d.Ticket.IDSolicitante,
		IDTecnicoAsignado:     d.Ticket.IDTecnicoAsignado,
		IDServicio:            d.Ticket.IDServicio,
		IDTipoTicket:          d.Ticket.IDTipoTicket,
		CodEstadoTicket:       d.Ticket.CodEstadoTicket,
		IDNivelPrioridad:      d.Ticket.IDNivelPrioridad,
		IDCatalogoFalla:       d.Ticket.IDCatalogoFalla,
		IDDepartamentoSoporte: d.Ticket.IDDepartamentoSoporte,
		Critico:               d.Ticket.Critico,
		DetalleFallaReportada: d.Ticket.DetalleFallaReportada,
		UbicacionObs:          d.Ticket.UbicacionObs,
		CreatedAt:             d.Ticket.CreatedAt.Format(timeFmt),
		UpdatedAt:             d.Ticket.UpdatedAt.Format(timeFmt),
		Trazabilidad:          trazabilidad,
		Bitacora:              bitacora,
	}

	if d.Ticket.FechaInicioTrabajo != nil {
		s := d.Ticket.FechaInicioTrabajo.Format(timeFmt)
		resp.FechaInicioTrabajo = &s
	}
	if d.Ticket.FechaFinTrabajo != nil {
		s := d.Ticket.FechaFinTrabajo.Format(timeFmt)
		resp.FechaFinTrabajo = &s
	}

	return resp
}

package models

import "time"

type EstadoTicket struct {
	ID              int    `gorm:"column:id;primaryKey"`
	Descripcion     string `gorm:"column:descripcion;not null;uniqueIndex"`
	CodEstadoTicket string `gorm:"column:cod_estado_ticket;not null;uniqueIndex"`
}

func (EstadoTicket) TableName() string { return "estado_ticket" }

type Ticket struct {
	ID                    int        `gorm:"column:id;primaryKey"`
	NroTicket             string     `gorm:"column:nro_ticket;not null;uniqueIndex"`
	IDSolicitante         int        `gorm:"column:id_solicitante;not null"`
	IDTecnicoAsignado     *int       `gorm:"column:id_tecnico_asignado"`
	IDServicio            *int       `gorm:"column:id_servicio"`
	IDTipoTicket          int        `gorm:"column:id_tipo_ticket;not null"`
	CodEstadoTicket       string     `gorm:"column:cod_estado_ticket;not null"`
	IDNivelPrioridad      *int       `gorm:"column:id_nivel_prioridad"`
	IDCatalogoFalla       *int       `gorm:"column:id_catalogo_falla"`
	IDDepartamentoSoporte *int       `gorm:"column:id_departamento_soporte"`
	Critico               bool       `gorm:"column:critico;default:false"`
	DetalleFallaReportada string     `gorm:"column:detalle_falla_reportada;not null"`
	UbicacionObs          string     `gorm:"column:ubicacion_obs;not null;default:SIN OBSERVACION"`
	CreatedAt             time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	FechaInicioTrabajo    *time.Time `gorm:"column:fecha_inicio_trabajo"`
	FechaFinTrabajo       *time.Time `gorm:"column:fecha_fin_trabajo"`
}

func (Ticket) TableName() string { return "ticket" }

type TrazabilidadTicket struct {
	ID                int       `gorm:"column:id;primaryKey"`
	IDTicket          int       `gorm:"column:id_ticket;not null"`
	CodEstadoTicket   string    `gorm:"column:cod_estado_ticket;not null"`
	RutResponsable    string    `gorm:"column:rut_responsable;not null"`
	FechaTrazabilidad time.Time `gorm:"column:fecha_trazabilidad;autoCreateTime"`
}

func (TrazabilidadTicket) TableName() string { return "trazabilidad_ticket" }

type BitacoraTicket struct {
	ID            int       `gorm:"column:id;primaryKey"`
	IDTicket      int       `gorm:"column:id_ticket;not null"`
	RutAutor      string    `gorm:"column:rut_autor;not null"`
	Comentario    string    `gorm:"column:comentario;not null"`
	FechaRegistro time.Time `gorm:"column:fecha_registro;autoCreateTime"`
}

func (BitacoraTicket) TableName() string { return "bitacora_ticket" }

type Valorizacion struct {
	ID            int       `gorm:"column:id;primaryKey"`
	IDTicket      int       `gorm:"column:id_ticket;not null;uniqueIndex"`
	IDTecnico     int       `gorm:"column:id_tecnico;not null"`
	IDSolicitante int       `gorm:"column:id_solicitante;not null"`
	Nota          int       `gorm:"column:nota;not null"`
	Comentarios   string    `gorm:"column:comentarios"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Valorizacion) TableName() string { return "valorizacion" }

type TicketPausa struct {
	ID                  int        `gorm:"column:id;primaryKey"`
	IDTicket            int        `gorm:"column:id_ticket;not null"`
	FechaInicioPausa    time.Time  `gorm:"column:fecha_inicio_pausa;autoCreateTime"`
	FechaFinPausa       *time.Time `gorm:"column:fecha_fin_pausa"`
	IDTecnicoPausa      int        `gorm:"column:id_tecnico_pausa;not null"`
	EstadoPausa         string     `gorm:"column:estado_pausa;not null;default:PENDIENTE"`
	IDMotivoPausa       int        `gorm:"column:id_motivo_pausa;not null"`
	IDTecnicoAutorizado *int       `gorm:"column:id_tecnico_autorizado"`
	FechaResolucion     *time.Time `gorm:"column:fecha_resolucion"`
}

func (TicketPausa) TableName() string { return "ticket_pausas" }

type TicketTraspaso struct {
	ID                   int        `gorm:"column:id;primaryKey"`
	IDTicket             int        `gorm:"column:id_ticket;not null"`
	IDTecnicoOrigen      int        `gorm:"column:id_tecnico_origen;not null"`
	IDTecnicoDestino     int        `gorm:"column:id_tecnico_destino;not null"`
	EstadoTraspaso       string     `gorm:"column:estado_traspaso;not null;default:PENDIENTE"`
	Motivo               string     `gorm:"column:motivo;not null"`
	ComentarioResolucion string     `gorm:"column:comentario_resolucion"`
	FechaSolicitud       time.Time  `gorm:"column:fecha_solicitud;autoCreateTime"`
	FechaResolucion      *time.Time `gorm:"column:fecha_resolucion"`
}

func (TicketTraspaso) TableName() string { return "ticket_traspasos" }

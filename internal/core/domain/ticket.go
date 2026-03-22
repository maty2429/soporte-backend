package domain

import "time"

type Ticket struct {
	ID                    int
	NroTicket             string
	IDSolicitante         int
	IDTecnicoAsignado     *int
	IDServicio            *int
	IDTipoTicket          int
	CodEstadoTicket       string
	IDNivelPrioridad      *int
	IDCatalogoFalla       *int
	IDDepartamentoSoporte *int
	Critico               bool
	DetalleFallaReportada string
	UbicacionObs          string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	FechaInicioTrabajo    *time.Time
	FechaFinTrabajo       *time.Time
}

type EstadoTicket struct {
	ID               int
	Descripcion      string
	CodEstadoTicket  string
}

type TrazabilidadTicket struct {
	ID                 int
	IDTicket           int
	CodEstadoTicket    string
	DescripcionEstado  string
	RutResponsable     string
	FechaTrazabilidad  time.Time
}

type BitacoraTicket struct {
	ID             int
	IDTicket       int
	RutAutor       string
	Comentario     string
	FechaRegistro  time.Time
}

type Valorizacion struct {
	ID             int
	IDTicket       int
	IDTecnico      int
	IDSolicitante  int
	Nota           int
	Comentarios    string
	CreatedAt      time.Time
}

type TicketPausa struct {
	ID                  int
	IDTicket            int
	FechaInicioPausa    time.Time
	FechaFinPausa       *time.Time
	IDTecnicoPausa      int
	EstadoPausa         string // PENDIENTE, APROBADA, RECHAZADA
	IDMotivoPausa       int
	IDTecnicoAutorizado *int
	FechaResolucion     *time.Time
}

type TicketTraspaso struct {
	ID                    int
	IDTicket              int
	IDTecnicoOrigen       int
	IDTecnicoDestino      int
	EstadoTraspaso        string // PENDIENTE, ACEPTADO, RECHAZADO
	Motivo                string
	ComentarioResolucion  string
	FechaSolicitud        time.Time
	FechaResolucion       *time.Time
}

type TicketDetalle struct {
	Ticket       Ticket
	Trazabilidad []TrazabilidadTicket
	Bitacora     []BitacoraTicket
}

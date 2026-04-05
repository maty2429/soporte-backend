package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

const (
	codEstadoCreado      = "CRE"
	codEstadoAsignado    = "ASI"
	codEstadoProgreso    = "PRO"
	codEstadoPausado     = "PAU"
	codEstadoTerminado   = "TER"
	codEstadoCerrado     = "CER"
	codEstadoSolTraspaso = "STR"
	codEstadoTraspasado  = "TRA"
	codEstadoCancelado   = "CAN"
	codEstadoVisto       = "VITEC"
	maxRetryNro          = 5
)

type TicketService struct {
	ticketRepo   ports.TicketRepository
	solRepo      ports.SolicitanteRepository
	catalogoRepo ports.CatalogoRepository
	tecnicoRepo  ports.TecnicoRepository
}

func NewTicketService(ticketRepo ports.TicketRepository, solRepo ports.SolicitanteRepository, catalogoRepo ports.CatalogoRepository, tecnicoRepo ports.TecnicoRepository) *TicketService {
	return &TicketService{
		ticketRepo:   ticketRepo,
		solRepo:      solRepo,
		catalogoRepo: catalogoRepo,
		tecnicoRepo:  tecnicoRepo,
	}
}

func (s *TicketService) Create(ctx context.Context, cmd CreateTicketCommand) (domain.Ticket, error) {
	if cmd.IDSolicitante <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_solicitante is required", nil)
	}
	if cmd.IDServicio <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_servicio is required", nil)
	}
	if cmd.IDTipoTicket <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_tipo_ticket is required", nil)
	}
	if cmd.IDNivelPrioridad <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_nivel_prioridad is required", nil)
	}
	if cmd.IDDepartamentoSoporte <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_departamento_soporte is required", nil)
	}
	detalle := strings.TrimSpace(cmd.DetalleFallaReportada)
	if detalle == "" {
		return domain.Ticket{}, domain.ValidationError("detalle_falla_reportada is required", nil)
	}

	sol, err := s.solRepo.GetByID(ctx, cmd.IDSolicitante)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get solicitante", err)
	}

	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoCreado)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get estado ticket", err)
	}

	ubicacionObs := "SIN OBSERVACION"
	if cmd.UbicacionObs != nil {
		trimmed := strings.TrimSpace(*cmd.UbicacionObs)
		if trimmed != "" {
			ubicacionObs = strings.ToUpper(trimmed)
		}
	}

	var ticket domain.Ticket
	for attempt := 0; attempt < maxRetryNro; attempt++ {
		nro := generateNroTicket()

		exists, err := s.ticketRepo.NroTicketExists(ctx, nro)
		if err != nil {
			return domain.Ticket{}, domain.InternalError("check nro_ticket", err)
		}
		if exists {
			continue
		}

		ticket = domain.Ticket{
			NroTicket:             nro,
			IDSolicitante:         cmd.IDSolicitante,
			IDServicio:            &cmd.IDServicio,
			IDTipoTicket:          cmd.IDTipoTicket,
			CodEstadoTicket:       estado.CodEstadoTicket,
			IDNivelPrioridad:      &cmd.IDNivelPrioridad,
			IDDepartamentoSoporte: &cmd.IDDepartamentoSoporte,
			Critico:               cmd.Critico,
			DetalleFallaReportada: detalle,
			UbicacionObs:          ubicacionObs,
		}

		txErr := s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
			if err := txRepo.Create(ctx, &ticket); err != nil {
				return err
			}

			traz := domain.TrazabilidadTicket{
				IDTicket:        ticket.ID,
				CodEstadoTicket: estado.CodEstadoTicket,
				RutResponsable:  sol.Rut + "-" + sol.Dv,
			}
			return txRepo.CreateTrazabilidad(ctx, &traz)
		})

		if txErr != nil {
			if isDuplicateKeyError(txErr) {
				continue
			}
			return domain.Ticket{}, wrapServiceError("create ticket", txErr)
		}

		return ticket, nil
	}

	return domain.Ticket{}, domain.InternalError("could not generate unique nro_ticket after retries", nil)
}

func (s *TicketService) GetByID(ctx context.Context, id int) (domain.Ticket, error) {
	if id <= 0 {
		return domain.Ticket{}, domain.ValidationError("id is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get ticket", err)
	}
	return ticket, nil
}

func (s *TicketService) GetByNroTicket(ctx context.Context, nroTicket string) (domain.TicketDetalle, error) {
	nro := strings.TrimSpace(strings.ToUpper(nroTicket))
	if nro == "" {
		return domain.TicketDetalle{}, domain.ValidationError("nro_ticket is required", nil)
	}

	ticket, err := s.ticketRepo.GetByNroTicket(ctx, nro)
	if err != nil {
		return domain.TicketDetalle{}, wrapServiceError("get ticket", err)
	}

	trazabilidad, err := s.ticketRepo.ListTrazabilidad(ctx, ticket.ID)
	if err != nil {
		return domain.TicketDetalle{}, wrapServiceError("list trazabilidad", err)
	}

	bitacora, err := s.ticketRepo.ListBitacora(ctx, ticket.ID)
	if err != nil {
		return domain.TicketDetalle{}, wrapServiceError("list bitacora", err)
	}

	return domain.TicketDetalle{
		Ticket:       ticket,
		Trazabilidad: trazabilidad,
		Bitacora:     bitacora,
	}, nil
}

func (s *TicketService) ListTickets(ctx context.Context, q ListTicketsQuery) (ListTicketsResult, error) {
	limit, offset := normalizePagination(q.Limit, q.Offset)

	rutTecnico := ""
	dvTecnico := ""
	if strings.TrimSpace(q.RutTecnico) != "" {
		var err error
		rutTecnico, dvTecnico, err = splitRutDV(q.RutTecnico)
		if err != nil {
			return ListTicketsResult{}, err
		}
	}

	items, total, err := s.ticketRepo.ListTickets(ctx, ports.ListTicketsFilters{
		CodEstadoTicket:       strings.TrimSpace(strings.ToUpper(q.CodEstadoTicket)),
		IDTecnicoAsignado:     q.IDTecnicoAsignado,
		RutTecnico:            rutTecnico,
		DVTecnico:             dvTecnico,
		IDSolicitante:         q.IDSolicitante,
		IDDepartamentoSoporte: q.IDDepartamentoSoporte,
		Critico:               q.Critico,
		Limit:                 limit,
		Offset:                offset,
	})
	if err != nil {
		return ListTicketsResult{}, wrapServiceError("list tickets", err)
	}

	return ListTicketsResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *TicketService) UpdateTicket(ctx context.Context, cmd UpdateTicketCommand) (domain.Ticket, error) {
	if cmd.IDTicket <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_ticket is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get ticket", err)
	}

	fields := make(map[string]any)

	if cmd.DetalleFallaReportada != nil {
		detalle := strings.TrimSpace(*cmd.DetalleFallaReportada)
		if detalle == "" {
			return domain.Ticket{}, domain.ValidationError("detalle_falla_reportada cannot be empty", nil)
		}
		fields["detalle_falla_reportada"] = detalle
		ticket.DetalleFallaReportada = detalle
	}

	if cmd.UbicacionObs != nil {
		ub := strings.TrimSpace(*cmd.UbicacionObs)
		if ub == "" {
			ub = "SIN OBSERVACION"
		} else {
			ub = strings.ToUpper(ub)
		}
		fields["ubicacion_obs"] = ub
		ticket.UbicacionObs = ub
	}

	if cmd.Critico != nil {
		fields["critico"] = *cmd.Critico
		ticket.Critico = *cmd.Critico
	}

	if cmd.IDTipoTicket != nil {
		if *cmd.IDTipoTicket <= 0 {
			return domain.Ticket{}, domain.ValidationError("id_tipo_ticket must be greater than 0", nil)
		}
		fields["id_tipo_ticket"] = *cmd.IDTipoTicket
		ticket.IDTipoTicket = *cmd.IDTipoTicket
	}

	if cmd.IDDepartamentoSoporte != nil {
		if *cmd.IDDepartamentoSoporte <= 0 {
			return domain.Ticket{}, domain.ValidationError("id_departamento_soporte must be greater than 0", nil)
		}
		fields["id_departamento_soporte"] = *cmd.IDDepartamentoSoporte
		ticket.IDDepartamentoSoporte = cmd.IDDepartamentoSoporte
	}

	if cmd.IDServicio != nil {
		if *cmd.IDServicio <= 0 {
			return domain.Ticket{}, domain.ValidationError("id_servicio must be greater than 0", nil)
		}
		fields["id_servicio"] = *cmd.IDServicio
		ticket.IDServicio = cmd.IDServicio
	}

	if len(fields) == 0 {
		return domain.Ticket{}, domain.ValidationError("no fields to update", nil)
	}

	if err := s.ticketRepo.UpdateFields(ctx, &ticket, fields); err != nil {
		return domain.Ticket{}, wrapServiceError("update ticket", err)
	}

	return ticket, nil
}

func (s *TicketService) Assign(ctx context.Context, cmd AssignTicketCommand) (domain.Ticket, error) {
	if cmd.IDTicket <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_ticket is required", nil)
	}
	if cmd.IDTecnicoAsignado <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_tecnico_asignado is required", nil)
	}
	if cmd.IDCatalogoFalla <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_catalogo_falla is required", nil)
	}
	if cmd.IDNivelPrioridad <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_nivel_prioridad is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get ticket", err)
	}

	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoAsignado)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get estado ticket", err)
	}

	sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
	if err != nil {
		return domain.Ticket{}, wrapServiceError("get solicitante", err)
	}

	ticket.IDTecnicoAsignado = &cmd.IDTecnicoAsignado
	ticket.IDCatalogoFalla = &cmd.IDCatalogoFalla
	ticket.IDNivelPrioridad = &cmd.IDNivelPrioridad
	ticket.CodEstadoTicket = estado.CodEstadoTicket

	txErr := s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		traz := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estado.CodEstadoTicket,
			RutResponsable:  sol.Rut + "-" + sol.Dv,
		}
		return txRepo.CreateTrazabilidad(ctx, &traz)
	})

	if txErr != nil {
		return domain.Ticket{}, wrapServiceError("assign ticket", txErr)
	}

	return ticket, nil
}

func (s *TicketService) ChangeEstado(ctx context.Context, cmd ChangeEstadoCommand) error {
	if cmd.IDTicket <= 0 {
		return domain.ValidationError("id_ticket is required", nil)
	}
	cod := strings.TrimSpace(strings.ToUpper(cmd.CodEstadoTicket))
	if cod == "" {
		return domain.ValidationError("cod_estado_ticket is required", nil)
	}
	rut := strings.TrimSpace(strings.ToUpper(cmd.RutResponsable))
	if rut == "" {
		return domain.ValidationError("rut_responsable is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return wrapServiceError("get ticket", err)
	}

	switch cod {
	case codEstadoVisto:
		if ticket.CodEstadoTicket != codEstadoAsignado {
			return domain.ValidationError("el ticket debe estar en ASIGNADO (ASI) para marcar como visto, estado actual: "+ticket.CodEstadoTicket, nil)
		}
	case codEstadoProgreso:
		if ticket.CodEstadoTicket != codEstadoVisto {
			return domain.ValidationError("el ticket debe estar en VISTO POR EL TÉCNICO (VITEC) para pasar a progreso, estado actual: "+ticket.CodEstadoTicket, nil)
		}
	case codEstadoCancelado:
		if ticket.CodEstadoTicket == codEstadoCancelado {
			return domain.ValidationError("el ticket ya está cancelado", nil)
		}
		if ticket.CodEstadoTicket == codEstadoCerrado {
			return domain.ValidationError("el ticket ya está cerrado, no se puede cancelar", nil)
		}
	case codEstadoTerminado:
		if ticket.CodEstadoTicket != codEstadoProgreso {
			return domain.ValidationError("el ticket debe estar en EN PROGRESO (PRO) para terminar, estado actual: "+ticket.CodEstadoTicket, nil)
		}
	default:
		return domain.ValidationError("cod_estado_ticket no permitido por este endpoint, valores válidos: VITEC, PRO, CAN, TER", nil)
	}

	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, cod)
	if err != nil {
		return wrapServiceError("get estado ticket", err)
	}

	ticket.CodEstadoTicket = estado.CodEstadoTicket

	return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		traz := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estado.CodEstadoTicket,
			RutResponsable:  rut,
		}
		return txRepo.CreateTrazabilidad(ctx, &traz)
	})
}

func (s *TicketService) Close(ctx context.Context, cmd CloseTicketCommand) error {
	if cmd.IDTicket <= 0 {
		return domain.ValidationError("id_ticket is required", nil)
	}
	if cmd.IDSolicitante <= 0 {
		return domain.ValidationError("id_solicitante is required", nil)
	}
	if cmd.Nota < 1 || cmd.Nota > 5 {
		return domain.ValidationError("nota must be between 1 and 5", nil)
	}
	observacion := strings.TrimSpace(cmd.Observacion)
	if observacion == "" {
		return domain.ValidationError("observacion is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return wrapServiceError("get ticket", err)
	}

	if ticket.CodEstadoTicket != codEstadoTerminado {
		return domain.ValidationError("el técnico no ha terminado su trabajo, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	if ticket.IDTecnicoAsignado == nil {
		return domain.ValidationError("ticket sin técnico asignado", nil)
	}

	sol, err := s.solRepo.GetByID(ctx, cmd.IDSolicitante)
	if err != nil {
		return wrapServiceError("get solicitante", err)
	}

	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoCerrado)
	if err != nil {
		return wrapServiceError("get estado ticket", err)
	}

	comentarios := strings.TrimSpace(cmd.Comentarios)
	rutResponsable := sol.Rut + "-" + sol.Dv

	return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		val := domain.Valorizacion{
			IDTicket:      ticket.ID,
			IDTecnico:     *ticket.IDTecnicoAsignado,
			IDSolicitante: cmd.IDSolicitante,
			Nota:          cmd.Nota,
			Comentarios:   comentarios,
		}
		if err := txRepo.CreateValorizacion(ctx, &val); err != nil {
			return err
		}

		bitacora := domain.BitacoraTicket{
			IDTicket:   ticket.ID,
			RutAutor:   strings.ToUpper(rutResponsable),
			Comentario: observacion,
		}
		if err := txRepo.CreateBitacora(ctx, &bitacora); err != nil {
			return err
		}

		ticket.CodEstadoTicket = estado.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		traz := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estado.CodEstadoTicket,
			RutResponsable:  rutResponsable,
		}
		return txRepo.CreateTrazabilidad(ctx, &traz)
	})
}

// generateNroTicket genera un número de ticket con formato TK-XXXXXX-YY
// donde XXXXXX son 6 dígitos calculados con nanosegundos + random
// y YY son los últimos 2 dígitos del año en curso.
func generateNroTicket() string {
	now := time.Now()
	nano := now.UnixNano()
	r := rand.Int63()
	num := ((nano ^ r) % 900000) + 100000
	if num < 0 {
		num = -num
	}
	num = (num % 900000) + 100000
	yy := now.Year() % 100
	return fmt.Sprintf("TK-%06d-%02d", num, yy)
}

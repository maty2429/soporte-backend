package services

import (
	"context"
	"errors"
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
	// --- validaciones ---
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

	// --- buscar solicitante (para obtener rut) ---
	sol, err := s.solRepo.GetByID(ctx, cmd.IDSolicitante)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Ticket{}, err
		}
		return domain.Ticket{}, domain.InternalError("get solicitante", err)
	}

	// --- buscar estado CRE ---
	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoCreado)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Ticket{}, err
		}
		return domain.Ticket{}, domain.InternalError("get estado ticket", err)
	}

	// --- ubicacion_obs ---
	ubicacionObs := "SIN OBSERVACION"
	if cmd.UbicacionObs != nil {
		trimmed := strings.TrimSpace(*cmd.UbicacionObs)
		if trimmed != "" {
			ubicacionObs = strings.ToUpper(trimmed)
		}
	}

	// --- generar nro_ticket con retry dentro de transacción ---
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
			var appErr *domain.Error
			if errors.As(txErr, &appErr) {
				return domain.Ticket{}, txErr
			}
			return domain.Ticket{}, domain.InternalError("create ticket", txErr)
		}

		return ticket, nil
	}

	return domain.Ticket{}, domain.InternalError("could not generate unique nro_ticket after retries", nil)
}

func (s *TicketService) GetByNroTicket(ctx context.Context, nroTicket string) (domain.TicketDetalle, error) {
	nro := strings.TrimSpace(strings.ToUpper(nroTicket))
	if nro == "" {
		return domain.TicketDetalle{}, domain.ValidationError("nro_ticket is required", nil)
	}

	ticket, err := s.ticketRepo.GetByNroTicket(ctx, nro)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.TicketDetalle{}, err
		}
		return domain.TicketDetalle{}, domain.InternalError("get ticket", err)
	}

	trazabilidad, err := s.ticketRepo.ListTrazabilidad(ctx, ticket.ID)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.TicketDetalle{}, err
		}
		return domain.TicketDetalle{}, domain.InternalError("list trazabilidad", err)
	}

	bitacora, err := s.ticketRepo.ListBitacora(ctx, ticket.ID)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.TicketDetalle{}, err
		}
		return domain.TicketDetalle{}, domain.InternalError("list bitacora", err)
	}

	return domain.TicketDetalle{
		Ticket:       ticket,
		Trazabilidad: trazabilidad,
		Bitacora:     bitacora,
	}, nil
}

func (s *TicketService) Assign(ctx context.Context, cmd AssignTicketCommand) (domain.Ticket, error) {
	// --- validaciones ---
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

	// --- buscar ticket ---
	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Ticket{}, err
		}
		return domain.Ticket{}, domain.InternalError("get ticket", err)
	}

	// --- buscar estado ASI ---
	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoAsignado)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Ticket{}, err
		}
		return domain.Ticket{}, domain.InternalError("get estado ticket", err)
	}

	// --- buscar solicitante (para rut en trazabilidad) ---
	sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.Ticket{}, err
		}
		return domain.Ticket{}, domain.InternalError("get solicitante", err)
	}

	// --- actualizar ticket + trazabilidad en transacción ---
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
		var appErr *domain.Error
		if errors.As(txErr, &appErr) {
			return domain.Ticket{}, txErr
		}
		return domain.Ticket{}, domain.InternalError("assign ticket", txErr)
	}

	return ticket, nil
}

func (s *TicketService) CreateBitacora(ctx context.Context, cmd CreateBitacoraCommand) (domain.BitacoraTicket, error) {
	if cmd.IDTicket <= 0 {
		return domain.BitacoraTicket{}, domain.ValidationError("id_ticket is required", nil)
	}
	rutAutor := strings.TrimSpace(cmd.RutAutor)
	if rutAutor == "" {
		return domain.BitacoraTicket{}, domain.ValidationError("rut_autor is required", nil)
	}
	comentario := strings.TrimSpace(cmd.Comentario)
	if comentario == "" {
		return domain.BitacoraTicket{}, domain.ValidationError("comentario is required", nil)
	}

	// verificar que el ticket existe
	if _, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket); err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.BitacoraTicket{}, err
		}
		return domain.BitacoraTicket{}, domain.InternalError("get ticket", err)
	}

	bitacora := domain.BitacoraTicket{
		IDTicket:   cmd.IDTicket,
		RutAutor:   strings.ToUpper(rutAutor),
		Comentario: comentario,
	}

	if err := s.ticketRepo.CreateBitacora(ctx, &bitacora); err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return domain.BitacoraTicket{}, err
		}
		return domain.BitacoraTicket{}, domain.InternalError("create bitacora", err)
	}

	return bitacora, nil
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

	// buscar ticket
	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return err
		}
		return domain.InternalError("get ticket", err)
	}

	// solo se permiten estos estados por este endpoint
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

	// buscar estado por código
	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, cod)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return err
		}
		return domain.InternalError("get estado ticket", err)
	}

	// update ticket + trazabilidad en transacción
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

func (s *TicketService) CreatePausa(ctx context.Context, cmd CreatePausaCommand) (domain.TicketPausa, error) {
	if cmd.IDTicket <= 0 {
		return domain.TicketPausa{}, domain.ValidationError("id_ticket is required", nil)
	}
	if cmd.IDTecnicoPausa <= 0 {
		return domain.TicketPausa{}, domain.ValidationError("id_tecnico_pausa is required", nil)
	}
	if cmd.IDMotivoPausa <= 0 {
		return domain.TicketPausa{}, domain.ValidationError("id_motivo_pausa is required", nil)
	}

	// buscar ticket
	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.TicketPausa{}, wrapDomainError("get ticket", err)
	}

	// verificar que esté en progreso
	if ticket.CodEstadoTicket != codEstadoProgreso {
		return domain.TicketPausa{}, domain.ValidationError("el ticket debe estar en progreso (PRO) para pausar, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	// verificar que el técnico sea el asignado
	if ticket.IDTecnicoAsignado == nil || *ticket.IDTecnicoAsignado != cmd.IDTecnicoPausa {
		return domain.TicketPausa{}, domain.ValidationError("el técnico no es el asignado al ticket", nil)
	}

	// buscar motivo de pausa para determinar si requiere autorización
	motivo, err := s.catalogoRepo.GetMotivoPausaByID(ctx, cmd.IDMotivoPausa)
	if err != nil {
		return domain.TicketPausa{}, wrapDomainError("get motivo pausa", err)
	}

	estadoPausa := "PENDIENTE"
	if !motivo.RequiereAutorizacion {
		estadoPausa = "APROBADA"
	}

	// buscar estado PAU
	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoPausado)
	if err != nil {
		return domain.TicketPausa{}, wrapDomainError("get estado ticket", err)
	}

	// buscar rut del técnico (solicitante repo no sirve, usamos el rut del técnico desde el ticket)
	sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
	if err != nil {
		return domain.TicketPausa{}, wrapDomainError("get solicitante", err)
	}

	pausa := domain.TicketPausa{
		IDTicket:       cmd.IDTicket,
		IDTecnicoPausa: cmd.IDTecnicoPausa,
		EstadoPausa:    estadoPausa,
		IDMotivoPausa:  cmd.IDMotivoPausa,
	}

	txErr := s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		if err := txRepo.CreatePausa(ctx, &pausa); err != nil {
			return err
		}

		ticket.CodEstadoTicket = estado.CodEstadoTicket
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
		return domain.TicketPausa{}, wrapDomainError("create pausa", txErr)
	}

	return pausa, nil
}

func (s *TicketService) ResolverPausa(ctx context.Context, cmd ResolverPausaCommand) error {
	if cmd.IDPausa <= 0 {
		return domain.ValidationError("id_pausa is required", nil)
	}
	if cmd.IDTecnicoAutorizado <= 0 {
		return domain.ValidationError("id_tecnico_autorizado is required", nil)
	}
	estado := strings.TrimSpace(strings.ToUpper(cmd.EstadoPausa))
	if estado != "APROBADA" && estado != "RECHAZADA" {
		return domain.ValidationError("estado_pausa must be APROBADA or RECHAZADA", nil)
	}

	// buscar pausa
	pausa, err := s.ticketRepo.GetPausaByID(ctx, cmd.IDPausa)
	if err != nil {
		return wrapDomainError("get pausa", err)
	}

	if pausa.EstadoPausa != "PENDIENTE" {
		return domain.ValidationError("la pausa ya fue resuelta: "+pausa.EstadoPausa, nil)
	}

	now := time.Now()
	pausa.EstadoPausa = estado
	pausa.IDTecnicoAutorizado = &cmd.IDTecnicoAutorizado
	pausa.FechaResolucion = &now

	// si es rechazada, el ticket vuelve a PRO
	if estado == "RECHAZADA" {
		ticket, err := s.ticketRepo.GetByID(ctx, pausa.IDTicket)
		if err != nil {
			return wrapDomainError("get ticket", err)
		}

		estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
		if err != nil {
			return wrapDomainError("get estado ticket", err)
		}

		sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
		if err != nil {
			return wrapDomainError("get solicitante", err)
		}

		return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
			if err := txRepo.UpdatePausa(ctx, &pausa); err != nil {
				return err
			}

			ticket.CodEstadoTicket = estadoPRO.CodEstadoTicket
			if err := txRepo.Update(ctx, &ticket); err != nil {
				return err
			}

			traz := domain.TrazabilidadTicket{
				IDTicket:        ticket.ID,
				CodEstadoTicket: estadoPRO.CodEstadoTicket,
				RutResponsable:  sol.Rut + "-" + sol.Dv,
			}
			return txRepo.CreateTrazabilidad(ctx, &traz)
		})
	}

	// si es aprobada, solo actualizar la pausa
	return s.ticketRepo.UpdatePausa(ctx, &pausa)
}

func (s *TicketService) ReanudarTicket(ctx context.Context, cmd ReanudarTicketCommand) error {
	if cmd.IDTicket <= 0 {
		return domain.ValidationError("id_ticket is required", nil)
	}
	if cmd.IDTecnicoPausa <= 0 {
		return domain.ValidationError("id_tecnico_pausa is required", nil)
	}

	// buscar ticket
	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return wrapDomainError("get ticket", err)
	}

	if ticket.CodEstadoTicket != codEstadoPausado {
		return domain.ValidationError("el ticket no está pausado, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	// buscar pausa activa (aprobada, sin fecha_fin)
	pausa, err := s.ticketRepo.GetPausaActiva(ctx, cmd.IDTicket)
	if err != nil {
		return wrapDomainError("get pausa activa", err)
	}

	// buscar estado PRO
	estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
	if err != nil {
		return wrapDomainError("get estado ticket", err)
	}

	sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
	if err != nil {
		return wrapDomainError("get solicitante", err)
	}

	now := time.Now()
	pausa.FechaFinPausa = &now

	return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		if err := txRepo.UpdatePausa(ctx, &pausa); err != nil {
			return err
		}

		ticket.CodEstadoTicket = estadoPRO.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		traz := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estadoPRO.CodEstadoTicket,
			RutResponsable:  sol.Rut + "-" + sol.Dv,
		}
		return txRepo.CreateTrazabilidad(ctx, &traz)
	})
}

func (s *TicketService) CreateTraspaso(ctx context.Context, cmd CreateTraspasoCommand) (domain.TicketTraspaso, error) {
	if cmd.IDTicket <= 0 {
		return domain.TicketTraspaso{}, domain.ValidationError("id_ticket is required", nil)
	}
	if cmd.IDTecnicoOrigen <= 0 {
		return domain.TicketTraspaso{}, domain.ValidationError("id_tecnico_origen is required", nil)
	}
	if cmd.IDTecnicoDestino <= 0 {
		return domain.TicketTraspaso{}, domain.ValidationError("id_tecnico_destino is required", nil)
	}
	if cmd.IDTecnicoOrigen == cmd.IDTecnicoDestino {
		return domain.TicketTraspaso{}, domain.ValidationError("id_tecnico_destino must be different from id_tecnico_origen", nil)
	}
	motivo := strings.TrimSpace(cmd.Motivo)
	if motivo == "" {
		return domain.TicketTraspaso{}, domain.ValidationError("motivo is required", nil)
	}

	// buscar ticket
	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.TicketTraspaso{}, wrapDomainError("get ticket", err)
	}

	// verificar que esté en PRO
	if ticket.CodEstadoTicket != codEstadoProgreso {
		return domain.TicketTraspaso{}, domain.ValidationError("el ticket debe estar en progreso (PRO) para solicitar traspaso, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	// verificar que el técnico origen sea el asignado
	if ticket.IDTecnicoAsignado == nil || *ticket.IDTecnicoAsignado != cmd.IDTecnicoOrigen {
		return domain.TicketTraspaso{}, domain.ValidationError("el técnico origen no es el asignado al ticket", nil)
	}

	// verificar que no haya traspaso pendiente
	_, err = s.ticketRepo.GetTraspasoPendiente(ctx, cmd.IDTicket)
	if err == nil {
		return domain.TicketTraspaso{}, domain.ValidationError("ya existe un traspaso pendiente para este ticket", nil)
	}

	// buscar estado STR
	estadoSTR, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoSolTraspaso)
	if err != nil {
		return domain.TicketTraspaso{}, wrapDomainError("get estado STR", err)
	}

	// buscar técnico origen para rut en trazabilidad
	tecnicoOrigen, err := s.tecnicoRepo.GetByID(ctx, cmd.IDTecnicoOrigen)
	if err != nil {
		return domain.TicketTraspaso{}, wrapDomainError("get tecnico origen", err)
	}

	traspaso := domain.TicketTraspaso{
		IDTicket:         cmd.IDTicket,
		IDTecnicoOrigen:  cmd.IDTecnicoOrigen,
		IDTecnicoDestino: cmd.IDTecnicoDestino,
		EstadoTraspaso:   "PENDIENTE",
		Motivo:           motivo,
	}

	txErr := s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		if err := txRepo.CreateTraspaso(ctx, &traspaso); err != nil {
			return err
		}

		ticket.CodEstadoTicket = estadoSTR.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		traz := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estadoSTR.CodEstadoTicket,
			RutResponsable:  tecnicoOrigen.Rut + "-" + tecnicoOrigen.Dv,
		}
		return txRepo.CreateTrazabilidad(ctx, &traz)
	})

	if txErr != nil {
		return domain.TicketTraspaso{}, wrapDomainError("create traspaso", txErr)
	}

	return traspaso, nil
}

func (s *TicketService) ResolverTraspaso(ctx context.Context, cmd ResolverTraspasoCommand) error {
	if cmd.IDTraspaso <= 0 {
		return domain.ValidationError("id_traspaso is required", nil)
	}
	if cmd.IDTecnicoDestino <= 0 {
		return domain.ValidationError("id_tecnico_destino is required", nil)
	}
	estado := strings.TrimSpace(strings.ToUpper(cmd.EstadoTraspaso))
	if estado != "ACEPTADO" && estado != "RECHAZADO" {
		return domain.ValidationError("estado_traspaso must be ACEPTADO or RECHAZADO", nil)
	}

	// buscar traspaso
	traspaso, err := s.ticketRepo.GetTraspasoByID(ctx, cmd.IDTraspaso)
	if err != nil {
		return wrapDomainError("get traspaso", err)
	}

	if traspaso.EstadoTraspaso != "PENDIENTE" {
		return domain.ValidationError("el traspaso ya fue resuelto: "+traspaso.EstadoTraspaso, nil)
	}
	if traspaso.IDTecnicoDestino != cmd.IDTecnicoDestino {
		return domain.ValidationError("el id_tecnico_destino no coincide con el traspaso", nil)
	}

	// buscar ticket
	ticket, err := s.ticketRepo.GetByID(ctx, traspaso.IDTicket)
	if err != nil {
		return wrapDomainError("get ticket", err)
	}

	// buscar técnico destino para trazabilidad de aceptación
	tecnicoDestino, err := s.tecnicoRepo.GetByID(ctx, traspaso.IDTecnicoDestino)
	if err != nil {
		return wrapDomainError("get tecnico destino", err)
	}
	rutResponsableDestino := tecnicoDestino.Rut + "-" + tecnicoDestino.Dv

	now := time.Now()
	traspaso.EstadoTraspaso = estado
	traspaso.ComentarioResolucion = strings.TrimSpace(cmd.ComentarioResolucion)
	traspaso.FechaResolucion = &now

	if estado == "RECHAZADO" {
		// rechazado: ticket vuelve a PRO
		estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
		if err != nil {
			return wrapDomainError("get estado PRO", err)
		}

		return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
			if err := txRepo.UpdateTraspaso(ctx, &traspaso); err != nil {
				return err
			}

			ticket.CodEstadoTicket = estadoPRO.CodEstadoTicket
			if err := txRepo.Update(ctx, &ticket); err != nil {
				return err
			}

			traz := domain.TrazabilidadTicket{
				IDTicket:        ticket.ID,
				CodEstadoTicket: estadoPRO.CodEstadoTicket,
				RutResponsable:  rutResponsableDestino,
			}
			return txRepo.CreateTrazabilidad(ctx, &traz)
		})
	}

	// aceptado: cambiar técnico + TRA + PRO (4 operaciones en tx)
	estadoTRA, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoTraspasado)
	if err != nil {
		return wrapDomainError("get estado TRA", err)
	}

	estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
	if err != nil {
		return wrapDomainError("get estado PRO", err)
	}

	return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		// 1. actualizar traspaso
		if err := txRepo.UpdateTraspaso(ctx, &traspaso); err != nil {
			return err
		}

		// 2. cambiar técnico asignado + estado TRA
		ticket.IDTecnicoAsignado = &traspaso.IDTecnicoDestino
		ticket.CodEstadoTicket = estadoTRA.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		// 3. trazabilidad TRA
		trazTRA := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estadoTRA.CodEstadoTicket,
			RutResponsable:  rutResponsableDestino,
		}
		if err := txRepo.CreateTrazabilidad(ctx, &trazTRA); err != nil {
			return err
		}

		// 4. cambiar estado a PRO automáticamente
		ticket.CodEstadoTicket = estadoPRO.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		// 5. trazabilidad PRO
		trazPRO := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estadoPRO.CodEstadoTicket,
			RutResponsable:  rutResponsableDestino,
		}
		return txRepo.CreateTrazabilidad(ctx, &trazPRO)
	})
}

func (s *TicketService) ListTraspasos(ctx context.Context, q ListTraspasosQuery) (ListTraspasosResult, error) {
	if q.IDTecnicoDestino <= 0 {
		return ListTraspasosResult{}, domain.ValidationError("id_tecnico_destino is required", nil)
	}

	estado := strings.TrimSpace(strings.ToUpper(q.Estado))
	if estado != "" && estado != "PENDIENTE" && estado != "ACEPTADO" && estado != "RECHAZADO" {
		return ListTraspasosResult{}, domain.ValidationError("estado must be PENDIENTE, ACEPTADO or RECHAZADO", nil)
	}

	limit := q.Limit
	if limit <= 0 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}
	offset := q.Offset
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.ticketRepo.ListTraspasos(ctx, ports.ListTraspasosFilters{
		IDTecnicoDestino: q.IDTecnicoDestino,
		Estado:           estado,
		Limit:            limit,
		Offset:           offset,
	})
	if err != nil {
		return ListTraspasosResult{}, wrapDomainError("list traspasos", err)
	}

	return ListTraspasosResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *TicketService) ListPausas(ctx context.Context, q ListPausasQuery) (ListPausasResult, error) {
	if q.IDTicket <= 0 {
		return ListPausasResult{}, domain.ValidationError("id_ticket is required", nil)
	}

	// validar estado si viene
	estado := strings.TrimSpace(strings.ToUpper(q.Estado))
	if estado != "" && estado != "PENDIENTE" && estado != "APROBADA" && estado != "RECHAZADA" {
		return ListPausasResult{}, domain.ValidationError("estado must be PENDIENTE, APROBADA or RECHAZADA", nil)
	}

	// normalizar limit/offset
	limit := q.Limit
	if limit <= 0 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}
	offset := q.Offset
	if offset < 0 {
		offset = 0
	}

	// verificar que el ticket existe
	if _, err := s.ticketRepo.GetByID(ctx, q.IDTicket); err != nil {
		return ListPausasResult{}, wrapDomainError("get ticket", err)
	}

	items, total, err := s.ticketRepo.ListPausas(ctx, ports.ListPausasFilters{
		IDTicket: q.IDTicket,
		Estado:   estado,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return ListPausasResult{}, wrapDomainError("list pausas", err)
	}

	return ListPausasResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// wrapDomainError helper para envolver errores del dominio

func (s *TicketService) UpdateTicket(ctx context.Context, cmd UpdateTicketCommand) (domain.Ticket, error) {
	if cmd.IDTicket <= 0 {
		return domain.Ticket{}, domain.ValidationError("id_ticket is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.Ticket{}, wrapDomainError("get ticket", err)
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
		return domain.Ticket{}, wrapDomainError("update ticket", err)
	}

	return ticket, nil
}

func (s *TicketService) GetByID(ctx context.Context, id int) (domain.Ticket, error) {
	if id <= 0 {
		return domain.Ticket{}, domain.ValidationError("id is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Ticket{}, wrapDomainError("get ticket", err)
	}
	return ticket, nil
}

func (s *TicketService) ListTickets(ctx context.Context, q ListTicketsQuery) (ListTicketsResult, error) {
	limit := q.Limit
	if limit <= 0 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}
	offset := q.Offset
	if offset < 0 {
		offset = 0
	}

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
		return ListTicketsResult{}, wrapDomainError("list tickets", err)
	}

	return ListTicketsResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *TicketService) ListBitacora(ctx context.Context, idTicket int) ([]domain.BitacoraTicket, error) {
	if idTicket <= 0 {
		return nil, domain.ValidationError("id_ticket is required", nil)
	}

	// verificar que el ticket existe
	if _, err := s.ticketRepo.GetByID(ctx, idTicket); err != nil {
		return nil, wrapDomainError("get ticket", err)
	}

	items, err := s.ticketRepo.ListBitacora(ctx, idTicket)
	if err != nil {
		return nil, wrapDomainError("list bitacora", err)
	}
	return items, nil
}

func wrapDomainError(msg string, err error) error {
	var appErr *domain.Error
	if errors.As(err, &appErr) {
		return err
	}
	return domain.InternalError(msg, err)
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

	// buscar ticket
	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return err
		}
		return domain.InternalError("get ticket", err)
	}

	// verificar que el estado sea TER (trabajo terminado)
	if ticket.CodEstadoTicket != codEstadoTerminado {
		return domain.ValidationError("el técnico no ha terminado su trabajo, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	// verificar que tenga técnico asignado
	if ticket.IDTecnicoAsignado == nil {
		return domain.ValidationError("ticket sin técnico asignado", nil)
	}

	// buscar solicitante para rut en trazabilidad
	sol, err := s.solRepo.GetByID(ctx, cmd.IDSolicitante)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return err
		}
		return domain.InternalError("get solicitante", err)
	}

	// buscar estado CER
	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoCerrado)
	if err != nil {
		var appErr *domain.Error
		if errors.As(err, &appErr) {
			return err
		}
		return domain.InternalError("get estado ticket", err)
	}

	comentarios := strings.TrimSpace(cmd.Comentarios)

	rutResponsable := sol.Rut + "-" + sol.Dv

	// transacción: valorizacion + bitácora + update ticket + trazabilidad
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
	// Asegurar rango 100000-999999
	num = (num % 900000) + 100000
	yy := now.Year() % 100
	return fmt.Sprintf("TK-%06d-%02d", num, yy)
}

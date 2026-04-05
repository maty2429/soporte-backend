package services

import (
	"context"
	"strings"
	"time"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

func (s *TicketService) ListTraspasos(ctx context.Context, q ListTraspasosQuery) (ListTraspasosResult, error) {
	if q.IDTecnicoDestino <= 0 {
		return ListTraspasosResult{}, domain.ValidationError("id_tecnico_destino is required", nil)
	}

	estado := strings.TrimSpace(strings.ToUpper(q.Estado))
	if estado != "" && estado != "PENDIENTE" && estado != "ACEPTADO" && estado != "RECHAZADO" {
		return ListTraspasosResult{}, domain.ValidationError("estado must be PENDIENTE, ACEPTADO or RECHAZADO", nil)
	}

	limit, offset := normalizePagination(q.Limit, q.Offset)

	items, total, err := s.ticketRepo.ListTraspasos(ctx, ports.ListTraspasosFilters{
		IDTecnicoDestino: q.IDTecnicoDestino,
		Estado:           estado,
		Limit:            limit,
		Offset:           offset,
	})
	if err != nil {
		return ListTraspasosResult{}, wrapServiceError("list traspasos", err)
	}

	return ListTraspasosResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
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

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.TicketTraspaso{}, wrapServiceError("get ticket", err)
	}

	if ticket.CodEstadoTicket != codEstadoProgreso {
		return domain.TicketTraspaso{}, domain.ValidationError("el ticket debe estar en progreso (PRO) para solicitar traspaso, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	if ticket.IDTecnicoAsignado == nil || *ticket.IDTecnicoAsignado != cmd.IDTecnicoOrigen {
		return domain.TicketTraspaso{}, domain.ValidationError("el técnico origen no es el asignado al ticket", nil)
	}

	_, err = s.ticketRepo.GetTraspasoPendiente(ctx, cmd.IDTicket)
	if err == nil {
		return domain.TicketTraspaso{}, domain.ValidationError("ya existe un traspaso pendiente para este ticket", nil)
	}

	estadoSTR, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoSolTraspaso)
	if err != nil {
		return domain.TicketTraspaso{}, wrapServiceError("get estado STR", err)
	}

	tecnicoOrigen, err := s.tecnicoRepo.GetByID(ctx, cmd.IDTecnicoOrigen)
	if err != nil {
		return domain.TicketTraspaso{}, wrapServiceError("get tecnico origen", err)
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
		return domain.TicketTraspaso{}, wrapServiceError("create traspaso", txErr)
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

	traspaso, err := s.ticketRepo.GetTraspasoByID(ctx, cmd.IDTraspaso)
	if err != nil {
		return wrapServiceError("get traspaso", err)
	}

	if traspaso.EstadoTraspaso != "PENDIENTE" {
		return domain.ValidationError("el traspaso ya fue resuelto: "+traspaso.EstadoTraspaso, nil)
	}
	if traspaso.IDTecnicoDestino != cmd.IDTecnicoDestino {
		return domain.ValidationError("el id_tecnico_destino no coincide con el traspaso", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, traspaso.IDTicket)
	if err != nil {
		return wrapServiceError("get ticket", err)
	}

	tecnicoDestino, err := s.tecnicoRepo.GetByID(ctx, traspaso.IDTecnicoDestino)
	if err != nil {
		return wrapServiceError("get tecnico destino", err)
	}
	rutResponsableDestino := tecnicoDestino.Rut + "-" + tecnicoDestino.Dv

	now := time.Now()
	traspaso.EstadoTraspaso = estado
	traspaso.ComentarioResolucion = strings.TrimSpace(cmd.ComentarioResolucion)
	traspaso.FechaResolucion = &now

	if estado == "RECHAZADO" {
		estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
		if err != nil {
			return wrapServiceError("get estado PRO", err)
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

	// aceptado: cambiar técnico + TRA + PRO en una transacción
	estadoTRA, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoTraspasado)
	if err != nil {
		return wrapServiceError("get estado TRA", err)
	}

	estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
	if err != nil {
		return wrapServiceError("get estado PRO", err)
	}

	return s.ticketRepo.RunInTx(ctx, func(txRepo ports.TicketRepository) error {
		if err := txRepo.UpdateTraspaso(ctx, &traspaso); err != nil {
			return err
		}

		ticket.IDTecnicoAsignado = &traspaso.IDTecnicoDestino
		ticket.CodEstadoTicket = estadoTRA.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		trazTRA := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estadoTRA.CodEstadoTicket,
			RutResponsable:  rutResponsableDestino,
		}
		if err := txRepo.CreateTrazabilidad(ctx, &trazTRA); err != nil {
			return err
		}

		ticket.CodEstadoTicket = estadoPRO.CodEstadoTicket
		if err := txRepo.Update(ctx, &ticket); err != nil {
			return err
		}

		trazPRO := domain.TrazabilidadTicket{
			IDTicket:        ticket.ID,
			CodEstadoTicket: estadoPRO.CodEstadoTicket,
			RutResponsable:  rutResponsableDestino,
		}
		return txRepo.CreateTrazabilidad(ctx, &trazPRO)
	})
}

package services

import (
	"context"
	"strings"
	"time"

	"soporte/internal/core/domain"
	"soporte/internal/core/ports"
)

func (s *TicketService) ListPausas(ctx context.Context, q ListPausasQuery) (ListPausasResult, error) {
	if q.IDTicket <= 0 {
		return ListPausasResult{}, domain.ValidationError("id_ticket is required", nil)
	}

	estado := strings.TrimSpace(strings.ToUpper(q.Estado))
	if estado != "" && estado != "PENDIENTE" && estado != "APROBADA" && estado != "RECHAZADA" {
		return ListPausasResult{}, domain.ValidationError("estado must be PENDIENTE, APROBADA or RECHAZADA", nil)
	}

	limit, offset := normalizePagination(q.Limit, q.Offset)

	if _, err := s.ticketRepo.GetByID(ctx, q.IDTicket); err != nil {
		return ListPausasResult{}, wrapServiceError("get ticket", err)
	}

	items, total, err := s.ticketRepo.ListPausas(ctx, ports.ListPausasFilters{
		IDTicket: q.IDTicket,
		Estado:   estado,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return ListPausasResult{}, wrapServiceError("list pausas", err)
	}

	return ListPausasResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
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

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return domain.TicketPausa{}, wrapServiceError("get ticket", err)
	}

	if ticket.CodEstadoTicket != codEstadoProgreso {
		return domain.TicketPausa{}, domain.ValidationError("el ticket debe estar en progreso (PRO) para pausar, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	if ticket.IDTecnicoAsignado == nil || *ticket.IDTecnicoAsignado != cmd.IDTecnicoPausa {
		return domain.TicketPausa{}, domain.ValidationError("el técnico no es el asignado al ticket", nil)
	}

	motivo, err := s.catalogoRepo.GetMotivoPausaByID(ctx, cmd.IDMotivoPausa)
	if err != nil {
		return domain.TicketPausa{}, wrapServiceError("get motivo pausa", err)
	}

	estadoPausa := "PENDIENTE"
	if !motivo.RequiereAutorizacion {
		estadoPausa = "APROBADA"
	}

	estado, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoPausado)
	if err != nil {
		return domain.TicketPausa{}, wrapServiceError("get estado ticket", err)
	}

	sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
	if err != nil {
		return domain.TicketPausa{}, wrapServiceError("get solicitante", err)
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
		return domain.TicketPausa{}, wrapServiceError("create pausa", txErr)
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

	pausa, err := s.ticketRepo.GetPausaByID(ctx, cmd.IDPausa)
	if err != nil {
		return wrapServiceError("get pausa", err)
	}

	if pausa.EstadoPausa != "PENDIENTE" {
		return domain.ValidationError("la pausa ya fue resuelta: "+pausa.EstadoPausa, nil)
	}

	now := time.Now()
	pausa.EstadoPausa = estado
	pausa.IDTecnicoAutorizado = &cmd.IDTecnicoAutorizado
	pausa.FechaResolucion = &now

	if estado == "RECHAZADA" {
		ticket, err := s.ticketRepo.GetByID(ctx, pausa.IDTicket)
		if err != nil {
			return wrapServiceError("get ticket", err)
		}

		estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
		if err != nil {
			return wrapServiceError("get estado ticket", err)
		}

		sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
		if err != nil {
			return wrapServiceError("get solicitante", err)
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

	return s.ticketRepo.UpdatePausa(ctx, &pausa)
}

func (s *TicketService) ReanudarTicket(ctx context.Context, cmd ReanudarTicketCommand) error {
	if cmd.IDTicket <= 0 {
		return domain.ValidationError("id_ticket is required", nil)
	}
	if cmd.IDTecnicoPausa <= 0 {
		return domain.ValidationError("id_tecnico_pausa is required", nil)
	}

	ticket, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket)
	if err != nil {
		return wrapServiceError("get ticket", err)
	}

	if ticket.CodEstadoTicket != codEstadoPausado {
		return domain.ValidationError("el ticket no está pausado, estado actual: "+ticket.CodEstadoTicket, nil)
	}

	pausa, err := s.ticketRepo.GetPausaActiva(ctx, cmd.IDTicket)
	if err != nil {
		return wrapServiceError("get pausa activa", err)
	}

	estadoPRO, err := s.ticketRepo.GetEstadoTicketByCod(ctx, codEstadoProgreso)
	if err != nil {
		return wrapServiceError("get estado ticket", err)
	}

	sol, err := s.solRepo.GetByID(ctx, ticket.IDSolicitante)
	if err != nil {
		return wrapServiceError("get solicitante", err)
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

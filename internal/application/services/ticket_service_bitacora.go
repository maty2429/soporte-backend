package services

import (
	"context"
	"strings"

	"soporte/internal/core/domain"
)

func (s *TicketService) ListBitacora(ctx context.Context, idTicket int) ([]domain.BitacoraTicket, error) {
	if idTicket <= 0 {
		return nil, domain.ValidationError("id_ticket is required", nil)
	}

	if _, err := s.ticketRepo.GetByID(ctx, idTicket); err != nil {
		return nil, wrapServiceError("get ticket", err)
	}

	items, err := s.ticketRepo.ListBitacora(ctx, idTicket)
	if err != nil {
		return nil, wrapServiceError("list bitacora", err)
	}
	return items, nil
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

	if _, err := s.ticketRepo.GetByID(ctx, cmd.IDTicket); err != nil {
		return domain.BitacoraTicket{}, wrapServiceError("get ticket", err)
	}

	bitacora := domain.BitacoraTicket{
		IDTicket:   cmd.IDTicket,
		RutAutor:   strings.ToUpper(rutAutor),
		Comentario: comentario,
	}

	if err := s.ticketRepo.CreateBitacora(ctx, &bitacora); err != nil {
		return domain.BitacoraTicket{}, wrapServiceError("create bitacora", err)
	}

	return bitacora, nil
}

package types

import "soporte/internal/core/domain"

const (
	DefaultListLimit = 20
	MaxListLimit     = 100
)

type ListSolicitantesQuery struct {
	Limit  int
	Offset int
	Search string
	Estado *bool
}

type CreateSolicitanteCommand struct {
	IDServicio     *int
	Correo         string
	Rut            string
	Dv             string
	NombreCompleto string
	Anexo          *int
	Estado         *bool
}

type UpdateSolicitanteCommand struct {
	ID             int
	IDServicio     *int
	Correo         *string
	Rut            *string
	Dv             *string
	NombreCompleto *string
	Anexo          *int
	Estado         *bool
}

type ListSolicitantesResult struct {
	Items  []domain.Solicitante
	Total  int64
	Limit  int
	Offset int
}

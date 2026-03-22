package services

import svctypes "soporte/internal/application/services/types"

const (
	DefaultListLimit = svctypes.DefaultListLimit
	MaxListLimit     = svctypes.MaxListLimit
)

type ListSolicitantesQuery = svctypes.ListSolicitantesQuery
type CreateSolicitanteCommand = svctypes.CreateSolicitanteCommand
type UpdateSolicitanteCommand = svctypes.UpdateSolicitanteCommand
type ListSolicitantesResult = svctypes.ListSolicitantesResult
type CreateTicketCommand = svctypes.CreateTicketCommand
type AssignTicketCommand = svctypes.AssignTicketCommand
type CreateBitacoraCommand = svctypes.CreateBitacoraCommand
type ChangeEstadoCommand = svctypes.ChangeEstadoCommand
type CreatePausaCommand = svctypes.CreatePausaCommand
type ResolverPausaCommand = svctypes.ResolverPausaCommand
type ReanudarTicketCommand = svctypes.ReanudarTicketCommand
type CloseTicketCommand = svctypes.CloseTicketCommand
type ListPausasQuery = svctypes.ListPausasQuery
type ListPausasResult = svctypes.ListPausasResult
type UpdateTicketCommand = svctypes.UpdateTicketCommand
type ListTicketsQuery = svctypes.ListTicketsQuery
type ListTicketsResult = svctypes.ListTicketsResult
type CreateTraspasoCommand = svctypes.CreateTraspasoCommand
type ResolverTraspasoCommand = svctypes.ResolverTraspasoCommand
type ListTraspasosQuery = svctypes.ListTraspasosQuery
type ListTraspasosResult = svctypes.ListTraspasosResult
type ListTecnicosQuery = svctypes.ListTecnicosQuery
type ListTecnicosResult = svctypes.ListTecnicosResult
type CreateTecnicoCommand = svctypes.CreateTecnicoCommand
type UpdateTecnicoCommand = svctypes.UpdateTecnicoCommand
type CreateHorarioTurnoCommand = svctypes.CreateHorarioTurnoCommand
type UpdateHorarioTurnoCommand = svctypes.UpdateHorarioTurnoCommand
type ListServiciosQuery = svctypes.ListServiciosQuery
type ListServiciosResult = svctypes.ListServiciosResult
type CreateServicioCommand = svctypes.CreateServicioCommand
type UpdateServicioCommand = svctypes.UpdateServicioCommand
type CreateCatalogoFallaCommand = svctypes.CreateCatalogoFallaCommand
type UpdateCatalogoFallaCommand = svctypes.UpdateCatalogoFallaCommand


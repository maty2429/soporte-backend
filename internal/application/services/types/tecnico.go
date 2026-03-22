package types

import "soporte/internal/core/domain"

type ListTecnicosQuery struct {
	Limit                 int
	Offset                int
	Search                string
	Estado                *bool
	IDTipoTecnico         int
	IDDepartamentoSoporte int
}

type ListTecnicosResult struct {
	Items  []domain.Tecnico
	Total  int64
	Limit  int
	Offset int
}

type CreateTecnicoCommand struct {
	Rut                   string
	Dv                    string
	NombreCompleto        string
	IDTipoTecnico         *int
	IDDepartamentoSoporte *int
	IDTipoTurno           *int
	Estado                *bool
}

type CreateHorarioTurnoCommand struct {
	IDTipoTurno int
	DiaSemana   int
	HoraInicio  string
	HoraFin     string
}

type UpdateHorarioTurnoCommand struct {
	ID          int
	IDTipoTurno *int
	DiaSemana   *int
	HoraInicio  *string
	HoraFin     *string
}

type UpdateTecnicoCommand struct {
	ID                    int
	Rut                   *string
	Dv                    *string
	NombreCompleto        *string
	IDTipoTecnico         *int
	IDDepartamentoSoporte *int
	IDTipoTurno           *int
	Estado                *bool
}

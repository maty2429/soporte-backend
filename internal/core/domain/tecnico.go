package domain

import "time"

type ConfiguracionHorarioTurno struct {
	ID           int
	IDTipoTurno  int
	DiaSemana    int
	HoraInicio   string
	HoraFin      string
}

type Tecnico struct {
	ID                    int
	Rut                   string
	Dv                    string
	NombreCompleto        string
	IDTipoTecnico         *int
	IDDepartamentoSoporte *int
	IDTipoTurno           *int
	Estado                bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

package domain

type ConfiguracionHorarioTurno struct {
	ID          int
	IDTipoTurno int
	DiaSemana   int
	HoraInicio  string
	HoraFin     string
}

type Tecnico struct {
	ID                    int
	Rut                   string
	Dv                    string
	NombreCompleto        string
	IDTipoTecnico         *int
	IDDepartamentoSoporte *int
	Estado                bool
	DepartamentoSoporte   *DepartamentoSoporte
}

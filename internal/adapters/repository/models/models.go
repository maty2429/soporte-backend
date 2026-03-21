package models

func All() []any {
	return []any{
		&NivelPrioridad{},
		&Servicio{},
		&Solicitante{},
		&TipoTicket{},
		&TipoTecnico{},
		&DepartamentoSoporte{},
		&MotivoPausa{},
		&TipoTurno{},
	}
}

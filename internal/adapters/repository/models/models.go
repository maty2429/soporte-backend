package models

func All() []any {
	return []any{
		&NivelPrioridad{},
		&Servicio{},
		&Solicitante{},
		&Tecnico{},
		&ConfiguracionHorarioTurno{},
		&TipoTicket{},
		&TipoTecnico{},
		&DepartamentoSoporte{},
		&MotivoPausa{},
		&CatalogoFalla{},
		&TipoTurno{},
		&EstadoTicket{},
		&Ticket{},
		&TrazabilidadTicket{},
		&BitacoraTicket{},
		&Valorizacion{},
		&TicketPausa{},
		&TicketTraspaso{},
	}
}

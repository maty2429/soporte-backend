package contracts

import (
	"testing"
	"time"

	"soporte/internal/core/domain"
)

func TestNewTicketResponseIncludesRelations(t *testing.T) {
	now := time.Date(2026, 3, 27, 10, 0, 0, 0, time.FixedZone("-03", -3*60*60))

	resp := NewTicketResponse(domain.Ticket{
		ID:                    1,
		NroTicket:             "TK-949247-26",
		IDSolicitante:         1,
		IDTecnicoAsignado:     intPtr(1),
		IDServicio:            intPtr(1),
		IDTipoTicket:          1,
		CodEstadoTicket:       "ASI",
		IDNivelPrioridad:      intPtr(1),
		IDCatalogoFalla:       intPtr(35),
		IDDepartamentoSoporte: intPtr(1),
		Critico:               false,
		DetalleFallaReportada: "No enciende el equipo",
		UbicacionObs:          "SALA DE SERVIDORES",
		CreatedAt:             now,
		UpdatedAt:             now,
		Solicitante: &domain.Solicitante{
			ID:             1,
			IDServicio:     intPtr(1),
			Correo:         "juan@example.com",
			Rut:            "12345678",
			Dv:             "K",
			NombreCompleto: "JUAN PEREZ",
			Estado:         true,
		},
		TecnicoAsignado: &domain.Tecnico{
			ID:             1,
			Rut:            "11111111",
			Dv:             "1",
			NombreCompleto: "TECNICO UNO",
			Estado:         true,
		},
		Servicio: &domain.Servicio{
			ID:        1,
			Edificio:  "PRINCIPAL",
			Piso:      1,
			Servicios: "SOPORTE",
			Ubicacion: "PISO 1",
			Unidades:  "TI",
		},
		TipoTicket: &domain.TipoTicket{
			ID:            1,
			CodTipoTicket: "INC",
			Descripcion:   "INCIDENCIA",
		},
		EstadoTicket: &domain.EstadoTicket{
			ID:              2,
			CodEstadoTicket: "ASI",
			Descripcion:     "ASIGNADO",
		},
		NivelPrioridad: &domain.NivelPrioridad{
			ID:          1,
			Descripcion: "ALTA",
		},
		CatalogoFalla: &domain.CatalogoFalla{
			ID:                   35,
			CodigoFalla:          "F035",
			DescripcionFalla:     "SIN ENERGIA",
			Complejidad:          3,
			RequiereVisitaFisica: true,
		},
		DepartamentoSoporte: &domain.DepartamentoSoporte{
			ID:              1,
			CodDepartamento: "TI",
			Descripcion:     "TECNOLOGIA",
		},
	})

	if resp.Solicitante == nil || resp.Solicitante.NombreCompleto != "JUAN PEREZ" {
		t.Fatalf("solicitante = %#v", resp.Solicitante)
	}
	if resp.TecnicoAsignado == nil || resp.TecnicoAsignado.NombreCompleto != "TECNICO UNO" {
		t.Fatalf("tecnico_asignado = %#v", resp.TecnicoAsignado)
	}
	if resp.Servicio == nil || resp.Servicio.Edificio != "PRINCIPAL" {
		t.Fatalf("servicio = %#v", resp.Servicio)
	}
	if resp.TipoTicket == nil || resp.TipoTicket.CodTipoTicket != "INC" {
		t.Fatalf("tipo_ticket = %#v", resp.TipoTicket)
	}
	if resp.EstadoTicket == nil || resp.EstadoTicket.CodEstadoTicket != "ASI" {
		t.Fatalf("estado_ticket = %#v", resp.EstadoTicket)
	}
	if resp.NivelPrioridad == nil || resp.NivelPrioridad.Descripcion != "ALTA" {
		t.Fatalf("nivel_prioridad = %#v", resp.NivelPrioridad)
	}
	if resp.CatalogoFalla == nil || resp.CatalogoFalla.CodigoFalla != "F035" {
		t.Fatalf("catalogo_falla = %#v", resp.CatalogoFalla)
	}
	if resp.DepartamentoSoporte == nil || resp.DepartamentoSoporte.CodDepartamento != "TI" {
		t.Fatalf("departamento_soporte = %#v", resp.DepartamentoSoporte)
	}
}

func intPtr(v int) *int {
	return &v
}

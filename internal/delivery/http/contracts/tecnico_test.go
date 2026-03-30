package contracts

import (
	"testing"
	"time"

	"soporte/internal/core/domain"
)

func TestNewTecnicoResponseIncludesDepartamentoSoporte(t *testing.T) {
	now := time.Date(2026, 3, 27, 10, 0, 0, 0, time.FixedZone("-03", -3*60*60))

	resp := NewTecnicoResponse(domain.Tecnico{
		ID:                    1,
		Rut:                   "11111111",
		Dv:                    "1",
		NombreCompleto:        "TECNICO UNO",
		IDDepartamentoSoporte: intPtr(1),
		Estado:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
		DepartamentoSoporte: &domain.DepartamentoSoporte{
			ID:              1,
			CodDepartamento: "TI",
			Descripcion:     "TECNOLOGIA",
		},
	})

	if resp.DepartamentoSoporte == nil {
		t.Fatal("departamento_soporte is nil")
	}
	if resp.DepartamentoSoporte.CodDepartamento != "TI" {
		t.Fatalf("cod_departamento = %q, want TI", resp.DepartamentoSoporte.CodDepartamento)
	}
}

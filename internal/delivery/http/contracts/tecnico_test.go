package contracts

import (
	"testing"

	"soporte/internal/core/domain"
)

func TestNewTecnicoResponseIncludesDepartamentoSoporte(t *testing.T) {
	resp := NewTecnicoResponse(domain.Tecnico{
		ID:                    1,
		Rut:                   "11111111",
		Dv:                    "1",
		NombreCompleto:        "TECNICO UNO",
		IDDepartamentoSoporte: intPtr(1),
		Estado:                true,
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

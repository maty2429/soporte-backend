package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	dbmodels "soporte/internal/adapters/repository/models"
	gormrepo "soporte/internal/adapters/repository/repository"
	"soporte/internal/application/services"
	"soporte/internal/core/domain"
)

func setupTecnicoService(t *testing.T) (*services.TecnicoService, *gorm.DB) {
	t.Helper()
	dbName := fmt.Sprintf("file:tecnico_%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.AutoMigrate(dbmodels.All()...); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// seed catalogs for FK references
	db.Exec("INSERT INTO tipo_tecnico (id, cod_tipo_tecnico, descripcion) VALUES (1, 'INF', 'INFORMATICO')")
	db.Exec("INSERT INTO departamento_soporte (id, cod_departamento_soporte, descripcion) VALUES (1, 'TI', 'TECNOLOGÍA')")
	db.Exec("INSERT INTO tipo_turno (id, cod_tipo_turno, descripcion) VALUES (1, 'DIA', 'DIURNO')")

	repo := gormrepo.NewTecnicoRepository(db)
	svc := services.NewTecnicoService(repo)
	return svc, db
}

func TestCreateTecnicoOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	tipoTec := 1
	depto := 1
	turno := 1

	tec, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:                   "11111111",
		Dv:                    "1",
		NombreCompleto:        "TECNICO UNO",
		IDTipoTecnico:         &tipoTec,
		IDDepartamentoSoporte: &depto,
		IDTipoTurno:           &turno,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if tec.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if tec.Rut != "11111111" {
		t.Errorf("rut = %q, want 11111111", tec.Rut)
	}
	if !tec.Estado {
		t.Error("estado should default to true")
	}
}

func TestCreateTecnicoDuplicateRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	cmd := services.CreateTecnicoCommand{
		Rut:            "22222222",
		Dv:             "2",
		NombreCompleto: "TECNICO DOS",
	}
	_, err := svc.Create(ctx, cmd)
	if err != nil {
		t.Fatalf("create 1: %v", err)
	}

	_, err = svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected conflict error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusConflict {
		t.Errorf("status = %d, want 409", appErr.Status())
	}
}

func TestCreateTecnicoMissingRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	_, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Dv:             "1",
		NombreCompleto: "SIN RUT",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestGetTecnicoByID(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	created, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "33333333",
		Dv:             "3",
		NombreCompleto: "TECNICO TRES",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.NombreCompleto != "TECNICO TRES" {
		t.Errorf("nombre = %q, want TECNICO TRES", got.NombreCompleto)
	}
}

func TestGetTecnicoByIDNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 9999)
	if err == nil {
		t.Fatal("expected not found")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestGetTecnicoByRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	_, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "44444444",
		Dv:             "4",
		NombreCompleto: "TECNICO CUATRO",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := svc.GetByRut(ctx, "44444444-4")
	if err != nil {
		t.Fatalf("get by rut: %v", err)
	}
	if got.NombreCompleto != "TECNICO CUATRO" {
		t.Errorf("nombre = %q, want TECNICO CUATRO", got.NombreCompleto)
	}
}

func TestUpdateTecnicoOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	created, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "55555555",
		Dv:             "5",
		NombreCompleto: "TECNICO CINCO",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	newName := "TECNICO CINCO ACTUALIZADO"
	updated, err := svc.Update(ctx, services.UpdateTecnicoCommand{
		ID:             created.ID,
		NombreCompleto: &newName,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.NombreCompleto != newName {
		t.Errorf("nombre = %q, want %q", updated.NombreCompleto, newName)
	}
}

func TestUpdateTecnicoDeactivate(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	created, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "66666666",
		Dv:             "6",
		NombreCompleto: "TECNICO SEIS",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	inactive := false
	updated, err := svc.Update(ctx, services.UpdateTecnicoCommand{
		ID:     created.ID,
		Estado: &inactive,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Estado {
		t.Error("expected estado = false")
	}
}

func TestUpdateTecnicoNoFields(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	created, err := svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "77777777",
		Dv:             "7",
		NombreCompleto: "TECNICO SIETE",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = svc.Update(ctx, services.UpdateTecnicoCommand{ID: created.ID})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestListTecnicosWithFilters(t *testing.T) {
	t.Parallel()
	svc, _ := setupTecnicoService(t)
	ctx := context.Background()

	tipoTec := 1
	// crear 2 técnicos
	_, _ = svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "88888888",
		Dv:             "8",
		NombreCompleto: "TECNICO OCHO",
		IDTipoTecnico:  &tipoTec,
	})
	_, _ = svc.Create(ctx, services.CreateTecnicoCommand{
		Rut:            "99999999",
		Dv:             "9",
		NombreCompleto: "TECNICO NUEVE",
	})

	// listar todos
	result, err := svc.List(ctx, services.ListTecnicosQuery{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("total = %d, want 2", result.Total)
	}

	// filtrar por tipo_tecnico
	result, err = svc.List(ctx, services.ListTecnicosQuery{IDTipoTecnico: 1})
	if err != nil {
		t.Fatalf("list filtered: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total filtered = %d, want 1", result.Total)
	}

	// search por nombre
	result, err = svc.List(ctx, services.ListTecnicosQuery{Search: "OCHO"})
	if err != nil {
		t.Fatalf("list search: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total search = %d, want 1", result.Total)
	}
}

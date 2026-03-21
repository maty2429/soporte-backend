package services_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"soporte/internal/adapters/repository/models"
	"soporte/internal/adapters/repository/repository"
	"soporte/internal/application/services"
	"soporte/internal/core/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSolicitanteService(t *testing.T) (*services.SolicitanteService, *gorm.DB) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", url.QueryEscape(t.Name()))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.AutoMigrate(models.All()...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}

	repo := repository.NewSolicitanteRepository(db)
	svc := services.NewSolicitanteService(repo)
	return svc, db
}

func boolPtr(b bool) *bool    { return &b }
func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }

func createServicio(t *testing.T, db *gorm.DB, nombre string) *int {
	t.Helper()

	row := models.Servicio{
		Edificio:  "Hospital",
		Piso:      2,
		Servicios: nombre,
		Ubicacion: "Torre A",
	}
	if err := db.Create(&row).Error; err != nil {
		t.Fatalf("create servicio: %v", err)
	}
	return &row.ID
}

// --- Create ---

func TestCreateSolicitante(t *testing.T) {
	t.Parallel()
	svc, db := setupSolicitanteService(t)
	idServicio := createServicio(t, db, "TI")

	sol, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "12345678",
		Dv:             "9",
		NombreCompleto: "Juan Perez",
		IDServicio:     idServicio,
		Correo:         "Juan@Test.com",
		Anexo:          intPtr(100),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sol.ID == 0 {
		t.Fatal("expected ID to be set")
	}
	if sol.Rut != "12345678" {
		t.Fatalf("expected rut=12345678, got %s", sol.Rut)
	}
	if sol.Correo != "juan@test.com" {
		t.Fatalf("expected correo lowercased, got %s", sol.Correo)
	}
	if !sol.Estado {
		t.Fatal("expected estado=true by default")
	}
	if sol.Anexo == nil || *sol.Anexo != 100 {
		t.Fatalf("expected anexo=100, got %v", sol.Anexo)
	}
	if sol.IDServicio == nil || *sol.IDServicio != *idServicio {
		t.Fatalf("expected id_servicio=%d, got %v", *idServicio, sol.IDServicio)
	}
}

func TestCreateSolicitanteEstadoFalse(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	sol, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "11111111",
		Dv:             "1",
		NombreCompleto: "Inactivo",
		Estado:         boolPtr(false),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sol.Estado {
		t.Fatal("expected estado=false")
	}
}

func TestCreateSolicitanteValidationRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "",
		Dv:             "9",
		NombreCompleto: "Test",
	})
	assertDomainError(t, err, http.StatusBadRequest)
}

func TestCreateSolicitanteValidationDv(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "12345678",
		Dv:             "",
		NombreCompleto: "Test",
	})
	assertDomainError(t, err, http.StatusBadRequest)
}

func TestCreateSolicitanteValidationNombreCompleto(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "12345678",
		Dv:             "9",
		NombreCompleto: "  ",
	})
	assertDomainError(t, err, http.StatusBadRequest)
}

func TestCreateSolicitanteDuplicateRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	cmd := services.CreateSolicitanteCommand{
		Rut:            "99999999",
		Dv:             "K",
		NombreCompleto: "Original",
	}
	if _, err := svc.Create(context.Background(), cmd); err != nil {
		t.Fatalf("first create: %v", err)
	}

	cmd.NombreCompleto = "Duplicado"
	_, err := svc.Create(context.Background(), cmd)
	assertDomainError(t, err, http.StatusConflict)
}

func TestCreateSolicitanteDuplicateCorreo(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "11111111",
		Dv:             "1",
		NombreCompleto: "Uno",
		Correo:         "dup@test.com",
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "22222222",
		Dv:             "2",
		NombreCompleto: "Dos",
		Correo:         "dup@test.com",
	})
	assertDomainError(t, err, http.StatusConflict)
}

// --- GetByID ---

func TestGetByID(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	created, _ := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "12345678",
		Dv:             "9",
		NombreCompleto: "Juan Perez",
	})

	sol, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sol.ID != created.ID {
		t.Fatalf("expected ID=%d, got %d", created.ID, sol.ID)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.GetByID(context.Background(), 9999)
	assertDomainError(t, err, http.StatusNotFound)
}

func TestGetByRut(t *testing.T) {
	t.Parallel()
	svc, db := setupSolicitanteService(t)
	idServicio := createServicio(t, db, "TI")

	created, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "12345678",
		Dv:             "9",
		NombreCompleto: "Juan Perez",
		IDServicio:     idServicio,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	sol, err := svc.GetByRut(context.Background(), "12.345.678-9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sol.ID != created.ID {
		t.Fatalf("expected ID=%d, got %d", created.ID, sol.ID)
	}
	if sol.Servicio == nil {
		t.Fatal("expected servicio relation loaded")
	}
	if sol.Servicio.Servicios != "TI" {
		t.Fatalf("expected servicio TI, got %s", sol.Servicio.Servicios)
	}
}

func TestGetByRutInvalidFormat(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.GetByRut(context.Background(), "12A")
	assertDomainError(t, err, http.StatusBadRequest)
}

// --- List ---

func TestListSolicitantes(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	for i := range 5 {
		rut := "1000000" + string(rune('0'+i))
		if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
			Rut:            rut,
			Dv:             "0",
			NombreCompleto: "Persona " + rut,
		}); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}

	result, err := svc.List(context.Background(), services.ListSolicitantesQuery{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 5 {
		t.Fatalf("expected total=5, got %d", result.Total)
	}
	if len(result.Items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(result.Items))
	}
	if result.Limit != services.DefaultListLimit {
		t.Fatalf("expected default limit=%d, got %d", services.DefaultListLimit, result.Limit)
	}
}

func TestListSolicitantesPagination(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	for i := range 5 {
		rut := "2000000" + string(rune('0'+i))
		if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
			Rut:            rut,
			Dv:             "0",
			NombreCompleto: "Persona " + rut,
		}); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}

	result, err := svc.List(context.Background(), services.ListSolicitantesQuery{
		Limit:  2,
		Offset: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 5 {
		t.Fatalf("expected total=5, got %d", result.Total)
	}
	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result.Items))
	}
}

func TestListSolicitantesSearch(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "11111111", Dv: "1", NombreCompleto: "Maria Lopez",
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "22222222", Dv: "2", NombreCompleto: "Pedro Garcia",
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	result, err := svc.List(context.Background(), services.ListSolicitantesQuery{
		Search: "maria",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("expected total=1, got %d", result.Total)
	}
	if result.Items[0].NombreCompleto != "Maria Lopez" {
		t.Fatalf("expected Maria Lopez, got %s", result.Items[0].NombreCompleto)
	}
}

func TestListSolicitantesFilterEstado(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "11111111", Dv: "1", NombreCompleto: "Activo",
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "22222222", Dv: "2", NombreCompleto: "Inactivo", Estado: boolPtr(false),
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	result, err := svc.List(context.Background(), services.ListSolicitantesQuery{
		Estado: boolPtr(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("expected total=1, got %d", result.Total)
	}
}

// --- Update ---

func TestUpdateSolicitante(t *testing.T) {
	t.Parallel()
	svc, db := setupSolicitanteService(t)
	idServicioInicial := createServicio(t, db, "TI")
	idServicioNuevo := createServicio(t, db, "RRHH")

	created, _ := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "12345678",
		Dv:             "9",
		NombreCompleto: "Original",
		IDServicio:     idServicioInicial,
	})

	updated, err := svc.Update(context.Background(), services.UpdateSolicitanteCommand{
		ID:             created.ID,
		NombreCompleto: strPtr("Actualizado"),
		IDServicio:     idServicioNuevo,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.NombreCompleto != "Actualizado" {
		t.Fatalf("expected NombreCompleto=Actualizado, got %s", updated.NombreCompleto)
	}
	if updated.IDServicio == nil || *updated.IDServicio != *idServicioNuevo {
		t.Fatalf("expected id_servicio=%d, got %v", *idServicioNuevo, updated.IDServicio)
	}
	// unchanged fields
	if updated.Rut != "12345678" {
		t.Fatalf("expected Rut unchanged, got %s", updated.Rut)
	}
}

func TestCreateSolicitanteInvalidServicio(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut:            "87654321",
		Dv:             "K",
		NombreCompleto: "Sin Servicio",
		IDServicio:     intPtr(9999),
	})
	assertDomainError(t, err, http.StatusBadRequest)
}

func TestUpdateSolicitanteNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	_, err := svc.Update(context.Background(), services.UpdateSolicitanteCommand{
		ID:             9999,
		NombreCompleto: strPtr("Nada"),
	})
	assertDomainError(t, err, http.StatusNotFound)
}

func TestUpdateSolicitanteNoFields(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	created, _ := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "12345678", Dv: "9", NombreCompleto: "Test",
	})

	_, err := svc.Update(context.Background(), services.UpdateSolicitanteCommand{
		ID: created.ID,
	})
	assertDomainError(t, err, http.StatusBadRequest)
}

func TestUpdateSolicitanteDeactivate(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	created, _ := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "12345678", Dv: "9", NombreCompleto: "Test",
	})

	updated, err := svc.Update(context.Background(), services.UpdateSolicitanteCommand{
		ID:     created.ID,
		Estado: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Estado {
		t.Fatal("expected estado=false after deactivation")
	}
}

func TestUpdateSolicitanteDuplicateRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupSolicitanteService(t)

	if _, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "11111111", Dv: "1", NombreCompleto: "Uno",
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	second, err := svc.Create(context.Background(), services.CreateSolicitanteCommand{
		Rut: "22222222", Dv: "2", NombreCompleto: "Dos",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = svc.Update(context.Background(), services.UpdateSolicitanteCommand{
		ID:  second.ID,
		Rut: strPtr("11111111"),
	})
	assertDomainError(t, err, http.StatusConflict)
}

// --- helpers ---

func assertDomainError(t *testing.T, err error, expectedStatus int) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var domainErr *domain.Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain.Error, got %T: %v", err, err)
	}
	if domainErr.Status() != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, domainErr.Status(), domainErr.Error())
	}
}

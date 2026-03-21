package services_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"soporte/internal/adapters/repository/models"
	"soporte/internal/adapters/repository/repository"
	"soporte/internal/application/services"
	"soporte/internal/core/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupCatalogoService(t *testing.T) *services.CatalogoService {
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

	repo := repository.NewCatalogoRepository(db)
	return services.NewCatalogoService(repo)
}

func assertCatalogoError(t *testing.T, err error, expectedStatus int) {
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

// ==================== niveles_prioridad ====================

func TestNivelesPrioridadCreate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	item, err := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{Descripcion: "Alta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == 0 {
		t.Fatal("expected ID to be set")
	}
	if item.Descripcion != "Alta" {
		t.Fatalf("expected Alta, got %v", item.Descripcion)
	}
}

func TestNivelesPrioridadList(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	for _, desc := range []string{"Alta", "Media", "Baja"} {
		if _, err := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{Descripcion: desc}); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	items, err := svc.ListNivelesPrioridad(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
}

func TestNivelesPrioridadCreateRequiredMissing(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	_, err := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{})
	assertCatalogoError(t, err, http.StatusBadRequest)
}

func TestNivelesPrioridadCreateDuplicate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	if _, err := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{Descripcion: "Alta"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{Descripcion: "Alta"})
	assertCatalogoError(t, err, http.StatusConflict)
}

func TestNivelesPrioridadUpdate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	created, _ := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{Descripcion: "Alta"})
	desc := "Muy Alta"
	updated, err := svc.UpdateNivelPrioridad(context.Background(), created.ID, &desc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Descripcion != "Muy Alta" {
		t.Fatalf("expected Muy Alta, got %v", updated.Descripcion)
	}
}

func TestNivelesPrioridadUpdateNotFound(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	desc := "Nada"
	_, err := svc.UpdateNivelPrioridad(context.Background(), 9999, &desc)
	assertCatalogoError(t, err, http.StatusNotFound)
}

func TestNivelesPrioridadUpdateNoFields(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	created, _ := svc.CreateNivelPrioridad(context.Background(), domain.NivelPrioridad{Descripcion: "Alta"})

	_, err := svc.UpdateNivelPrioridad(context.Background(), created.ID, nil)
	assertCatalogoError(t, err, http.StatusBadRequest)
}

// ==================== tipos_turno ====================

func TestTiposTurnoCreate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	item, err := svc.CreateTipoTurno(context.Background(), domain.TipoTurno{
		Nombre:      "Diurno",
		Descripcion: "Turno de día",
		Estado:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Nombre != "Diurno" {
		t.Fatalf("expected Diurno, got %v", item.Nombre)
	}
	if !item.Estado {
		t.Fatal("expected estado=true")
	}
}

func TestTiposTurnoCreateEstadoFalse(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	item, err := svc.CreateTipoTurno(context.Background(), domain.TipoTurno{
		Nombre: "Nocturno",
		Estado: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Estado {
		t.Fatal("expected estado=false")
	}
}

func TestTiposTurnoUpdate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	created, _ := svc.CreateTipoTurno(context.Background(), domain.TipoTurno{Nombre: "Diurno", Estado: true})

	estado := false
	updated, err := svc.UpdateTipoTurno(context.Background(), created.ID, nil, nil, &estado)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Estado {
		t.Fatal("expected estado=false")
	}
	if updated.Nombre != "Diurno" {
		t.Fatalf("expected nombre unchanged, got %v", updated.Nombre)
	}
}

func TestTiposTurnoMaxLen(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	longName := strings.Repeat("a", 51)
	_, err := svc.CreateTipoTurno(context.Background(), domain.TipoTurno{Nombre: longName, Estado: true})
	assertCatalogoError(t, err, http.StatusBadRequest)
}

// ==================== motivos_pausa ====================

func TestMotivosPausaCreate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	item, err := svc.CreateMotivoPausa(context.Background(), domain.MotivoPausa{
		MotivoPausa:          "Almuerzo",
		RequiereAutorizacion: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !item.RequiereAutorizacion {
		t.Fatal("expected requiere_autorizacion=true")
	}
}

// ==================== departamentos_soporte ====================

func TestDepartamentosSoporteCreate(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	item, err := svc.CreateDepartamentoSoporte(context.Background(), domain.DepartamentoSoporte{
		CodDepartamento: "INF",
		Descripcion:     "Informática",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.CodDepartamento != "INF" {
		t.Fatalf("expected INF, got %v", item.CodDepartamento)
	}
}

func TestDepartamentosSoporteCreatePartialRequired(t *testing.T) {
	t.Parallel()
	svc := setupCatalogoService(t)

	_, err := svc.CreateDepartamentoSoporte(context.Background(), domain.DepartamentoSoporte{CodDepartamento: "INF"})
	assertCatalogoError(t, err, http.StatusBadRequest)
}

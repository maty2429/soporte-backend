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

func setupTicketService(t *testing.T) (*services.TicketService, *gorm.DB) {
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

	// seed estado_ticket
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (1, 'CREADO', 'CRE')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (2, 'ASIGNADO', 'ASI')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (3, 'EN PROGRESO', 'PRO')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (4, 'PAUSADO', 'PAU')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (5, 'RESUELTO', 'RES')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (6, 'CERRADO', 'CER')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (7, 'CANCELADO', 'CAN')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (8, 'TRABAJO TERMINADO', 'TER')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (9, 'REABIERTO', 'REA')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (10, 'VISTO POR EL TÉCNICO', 'VITEC')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (11, 'TRASPASADO A OTRO TÉCNICO', 'TRA')")
	db.Exec("INSERT INTO estado_ticket (id, descripcion, cod_estado_ticket) VALUES (12, 'SOLICITUD DE TRASPASO', 'STR')")

	// seed tipo_ticket
	db.Exec("INSERT INTO tipo_ticket (id, cod_tipo_ticket, descripcion) VALUES (1, 'INC', 'INCIDENCIA')")

	// seed niveles_prioridad
	db.Exec("INSERT INTO niveles_prioridad (id, descripcion) VALUES (1, 'ALTA')")

	// seed departamentos_soporte
	db.Exec("INSERT INTO departamentos_soporte (id, cod_departamento, descripcion) VALUES (1, 'TI', 'TECNOLOGIA')")

	// seed servicio
	db.Exec("INSERT INTO servicio (id, edificio, piso, servicios, ubicacion) VALUES (1, 'PRINCIPAL', 1, 'SOPORTE', 'PISO 1')")

	// seed solicitante
	db.Exec("INSERT INTO solicitantes (id, id_servicio, rut, dv, nombre_completo, estado) VALUES (1, 1, '12345678', 'K', 'JUAN PEREZ', 1)")

	// tablas referenciadas por ticket (no tienen modelo GORM aún)
	db.Exec("CREATE TABLE IF NOT EXISTS tecnicos (id INTEGER PRIMARY KEY, rut TEXT, dv TEXT, nombre_completo TEXT)")
	db.Exec("INSERT INTO tecnicos (id, rut, dv, nombre_completo) VALUES (1, '11111111', '1', 'TECNICO UNO')")
	db.Exec("INSERT INTO tecnicos (id, rut, dv, nombre_completo) VALUES (2, '22222222', '2', 'TECNICO DOS')")
	db.Exec("CREATE TABLE IF NOT EXISTS catalogo_fallas (id INTEGER PRIMARY KEY, codigo_falla TEXT, descripcion_falla TEXT)")
	db.Exec("INSERT INTO catalogo_fallas (id, codigo_falla, descripcion_falla) VALUES (1, 'F001', 'FALLA DE RED')")

	// seed motivos_pausa
	db.Exec("INSERT INTO motivos_pausa (id, motivo_pausa, requiere_autorizacion) VALUES (1, 'Falta repuesto en bodega', 1)")
	db.Exec("INSERT INTO motivos_pausa (id, motivo_pausa, requiere_autorizacion) VALUES (2, 'Esperando al usuario', 0)")

	solRepo := repository.NewSolicitanteRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	catalogoRepo := repository.NewCatalogoRepository(db)
	svc := services.NewTicketService(ticketRepo, solRepo, catalogoRepo)
	return svc, db
}

func validTicketCmd() services.CreateTicketCommand {
	return services.CreateTicketCommand{
		IDSolicitante:         1,
		IDServicio:            1,
		IDTipoTicket:          1,
		IDNivelPrioridad:      1,
		IDDepartamentoSoporte: 1,
		Critico:               false,
		DetalleFallaReportada: "No enciende el equipo",
	}
}

// --- Tests ---

func TestTicketCreateOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket, err := svc.Create(ctx, validTicketCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verificar formato TK-XXXXXX-YY
	if len(ticket.NroTicket) != 12 {
		t.Errorf("nro_ticket length = %d, want 12, got %q", len(ticket.NroTicket), ticket.NroTicket)
	}
	if ticket.NroTicket[:3] != "TK-" {
		t.Errorf("nro_ticket should start with TK-, got %q", ticket.NroTicket)
	}
	if ticket.NroTicket[9] != '-' {
		t.Errorf("nro_ticket should have - at position 9, got %q", ticket.NroTicket)
	}
	if ticket.ID == 0 {
		t.Error("ticket ID should not be 0")
	}
	if ticket.CodEstadoTicket != "CRE" {
		t.Errorf("estado = %q, want CRE", ticket.CodEstadoTicket)
	}
}

func TestTicketCreateRegistraTrazabilidad(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket, err := svc.Create(ctx, validTicketCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verificar que existe trazabilidad
	var traz models.TrazabilidadTicket
	if err := db.Where("id_ticket = ?", ticket.ID).Take(&traz).Error; err != nil {
		t.Fatalf("trazabilidad not found: %v", err)
	}
	if traz.CodEstadoTicket != "CRE" {
		t.Errorf("trazabilidad estado = %q, want CRE", traz.CodEstadoTicket)
	}
	if traz.RutResponsable != "12345678-K" {
		t.Errorf("trazabilidad rut = %q, want 12345678-K", traz.RutResponsable)
	}
}

func TestTicketCreateUbicacionObs(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	obs := "Sala de servidores"
	cmd := validTicketCmd()
	cmd.UbicacionObs = &obs

	ticket, err := svc.Create(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var row models.Ticket
	db.Where("id = ?", ticket.ID).Take(&row)
	if row.UbicacionObs != "SALA DE SERVIDORES" {
		t.Errorf("ubicacion_obs = %q, want SALA DE SERVIDORES", row.UbicacionObs)
	}
}

func TestTicketCreateDefaultUbicacionObs(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket, err := svc.Create(ctx, validTicketCmd())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ticket.UbicacionObs != "SIN OBSERVACION" {
		t.Errorf("ubicacion_obs = %q, want SIN OBSERVACION", ticket.UbicacionObs)
	}
}

func TestTicketCreateCritico(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.Critico = true

	ticket, err := svc.Create(ctx, cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ticket.Critico {
		t.Error("critico should be true")
	}
}

func TestTicketCreateNroTicketUnique(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	// Crear varios tickets y verificar que todos tengan nro_ticket distinto
	nros := make(map[string]bool)
	for i := 0; i < 10; i++ {
		ticket, err := svc.Create(ctx, validTicketCmd())
		if err != nil {
			t.Fatalf("ticket %d: unexpected error: %v", i, err)
		}
		if nros[ticket.NroTicket] {
			t.Fatalf("duplicate nro_ticket: %s", ticket.NroTicket)
		}
		nros[ticket.NroTicket] = true
	}
}

func TestTicketCreateMissingSolicitante(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.IDSolicitante = 0

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketCreateMissingDetalle(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.DetalleFallaReportada = "   "

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketCreateSolicitanteNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.IDSolicitante = 999

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestTicketCreateMissingServicio(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.IDServicio = 0

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketCreateMissingTipoTicket(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.IDTipoTicket = 0

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketCreateMissingNivelPrioridad(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.IDNivelPrioridad = 0

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketCreateMissingDepartamento(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	cmd := validTicketCmd()
	cmd.IDDepartamentoSoporte = 0

	_, err := svc.Create(ctx, cmd)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ==================== Assign ====================

func createTicketForAssign(t *testing.T, svc *services.TicketService) domain.Ticket {
	t.Helper()
	ticket, err := svc.Create(context.Background(), validTicketCmd())
	if err != nil {
		t.Fatalf("create ticket: %v", err)
	}
	return ticket
}

func TestTicketAssignOK(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	result, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          ticket.ID,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.CodEstadoTicket != "ASI" {
		t.Errorf("estado = %q, want ASI", result.CodEstadoTicket)
	}
	if result.IDTecnicoAsignado == nil || *result.IDTecnicoAsignado != 1 {
		t.Errorf("id_tecnico_asignado = %v, want 1", result.IDTecnicoAsignado)
	}
	if result.IDCatalogoFalla == nil || *result.IDCatalogoFalla != 1 {
		t.Errorf("id_catalogo_falla = %v, want 1", result.IDCatalogoFalla)
	}
	if result.IDNivelPrioridad == nil || *result.IDNivelPrioridad != 1 {
		t.Errorf("id_nivel_prioridad = %v, want 1", result.IDNivelPrioridad)
	}

	// Verificar trazabilidad ASI
	var count int64
	db.Model(&models.TrazabilidadTicket{}).Where("id_ticket = ? AND cod_estado_ticket = ?", ticket.ID, "ASI").Count(&count)
	if count != 1 {
		t.Errorf("trazabilidad ASI count = %d, want 1", count)
	}
}

func TestTicketAssignRegistraTrazabilidad(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	_, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          ticket.ID,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var traz models.TrazabilidadTicket
	if err := db.Where("id_ticket = ? AND cod_estado_ticket = ?", ticket.ID, "ASI").Take(&traz).Error; err != nil {
		t.Fatalf("trazabilidad ASI not found: %v", err)
	}
	if traz.RutResponsable != "12345678-K" {
		t.Errorf("rut_responsable = %q, want 12345678-K", traz.RutResponsable)
	}
}

func TestTicketAssignTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          999,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestTicketAssignMissingTecnico(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          1,
		IDTecnicoAsignado: 0,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketAssignMissingCatalogoFalla(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          1,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   0,
		IDNivelPrioridad:  1,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestTicketAssignMissingNivelPrioridad(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          1,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  0,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ==================== Bitacora ====================

func TestBitacoraCreateOK(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	bitacora, err := svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   ticket.ID,
		RutAutor:   "12345678-K",
		Comentario: "Se revisó el equipo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bitacora.ID == 0 {
		t.Error("bitacora ID should not be 0")
	}
	if bitacora.IDTicket != ticket.ID {
		t.Errorf("id_ticket = %d, want %d", bitacora.IDTicket, ticket.ID)
	}
	if bitacora.RutAutor != "12345678-K" {
		t.Errorf("rut_autor = %q, want 12345678-K", bitacora.RutAutor)
	}
	if bitacora.Comentario != "Se revisó el equipo" {
		t.Errorf("comentario = %q, want 'Se revisó el equipo'", bitacora.Comentario)
	}

	// Verificar en BD
	var row models.BitacoraTicket
	if err := db.Where("id = ?", bitacora.ID).Take(&row).Error; err != nil {
		t.Fatalf("bitacora not found: %v", err)
	}
}

func TestBitacoraCreateTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   999,
		RutAutor:   "12345678-K",
		Comentario: "comentario",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestBitacoraCreateMissingComentario(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   1,
		RutAutor:   "12345678-K",
		Comentario: "   ",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestBitacoraCreateMissingRutAutor(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   1,
		RutAutor:   "",
		Comentario: "comentario",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestBitacoraCreateMissingTicket(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   0,
		RutAutor:   "12345678-K",
		Comentario: "comentario",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ==================== GetByNroTicket ====================

func TestGetByNroTicketOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	// crear ticket
	ticket, err := svc.Create(ctx, validTicketCmd())
	if err != nil {
		t.Fatalf("create ticket: %v", err)
	}

	// agregar bitacora
	_, err = svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   ticket.ID,
		RutAutor:   "12345678-K",
		Comentario: "Primer comentario",
	})
	if err != nil {
		t.Fatalf("create bitacora: %v", err)
	}

	// buscar por nro_ticket
	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if detalle.Ticket.NroTicket != ticket.NroTicket {
		t.Errorf("nro_ticket = %q, want %q", detalle.Ticket.NroTicket, ticket.NroTicket)
	}
	if detalle.Ticket.ID != ticket.ID {
		t.Errorf("id = %d, want %d", detalle.Ticket.ID, ticket.ID)
	}

	// debe tener 1 trazabilidad (CRE al crear)
	if len(detalle.Trazabilidad) != 1 {
		t.Errorf("trazabilidad count = %d, want 1", len(detalle.Trazabilidad))
	} else {
		if detalle.Trazabilidad[0].CodEstadoTicket != "CRE" {
			t.Errorf("trazabilidad estado = %q, want CRE", detalle.Trazabilidad[0].CodEstadoTicket)
		}
		if detalle.Trazabilidad[0].DescripcionEstado != "CREADO" {
			t.Errorf("trazabilidad descripcion = %q, want CREADO", detalle.Trazabilidad[0].DescripcionEstado)
		}
	}

	// debe tener 1 bitacora
	if len(detalle.Bitacora) != 1 {
		t.Errorf("bitacora count = %d, want 1", len(detalle.Bitacora))
	} else {
		if detalle.Bitacora[0].Comentario != "Primer comentario" {
			t.Errorf("bitacora comentario = %q, want 'Primer comentario'", detalle.Bitacora[0].Comentario)
		}
	}
}

func TestGetByNroTicketWithAssign(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket, err := svc.Create(ctx, validTicketCmd())
	if err != nil {
		t.Fatalf("create ticket: %v", err)
	}

	// asignar
	_, err = svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          ticket.ID,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err != nil {
		t.Fatalf("assign: %v", err)
	}

	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// debe tener 2 trazabilidades: CRE + ASI
	if len(detalle.Trazabilidad) != 2 {
		t.Errorf("trazabilidad count = %d, want 2", len(detalle.Trazabilidad))
	}
	if detalle.Ticket.CodEstadoTicket != "ASI" {
		t.Errorf("estado = %q, want ASI", detalle.Ticket.CodEstadoTicket)
	}
}

func TestGetByNroTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.GetByNroTicket(ctx, "TK-999999-26")
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestGetByNroTicketEmpty(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.GetByNroTicket(ctx, "  ")
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestGetByNroTicketEmptyBitacoraAndTrazabilidad(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket, err := svc.Create(ctx, validTicketCmd())
	if err != nil {
		t.Fatalf("create ticket: %v", err)
	}

	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// trazabilidad tiene 1 (CRE se crea al crear ticket)
	if len(detalle.Trazabilidad) != 1 {
		t.Errorf("trazabilidad count = %d, want 1", len(detalle.Trazabilidad))
	}
	// bitacora vacía
	if len(detalle.Bitacora) != 0 {
		t.Errorf("bitacora count = %d, want 0", len(detalle.Bitacora))
	}
}

// ==================== ChangeEstado ====================

func TestChangeEstadoOK(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "CAN",
		RutResponsable:  "12345678-K",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verificar que el ticket cambió de estado
	var row models.Ticket
	if err := db.Where("id = ?", ticket.ID).Take(&row).Error; err != nil {
		t.Fatalf("ticket not found: %v", err)
	}
	if row.CodEstadoTicket != "CAN" {
		t.Errorf("estado = %q, want CAN", row.CodEstadoTicket)
	}

	// Verificar trazabilidad CAN
	var traz models.TrazabilidadTicket
	if err := db.Where("id_ticket = ? AND cod_estado_ticket = ?", ticket.ID, "CAN").Take(&traz).Error; err != nil {
		t.Fatalf("trazabilidad CAN not found: %v", err)
	}
	if traz.RutResponsable != "12345678-K" {
		t.Errorf("rut_responsable = %q, want 12345678-K", traz.RutResponsable)
	}
}

func TestChangeEstadoTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        999,
		CodEstadoTicket: "PRO",
		RutResponsable:  "12345678-K",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestChangeEstadoInvalidCode(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "INVALIDO",
		RutResponsable:  "12345678-K",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestChangeEstadoMissingCod(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        1,
		CodEstadoTicket: "",
		RutResponsable:  "12345678-K",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestChangeEstadoMissingRut(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        1,
		CodEstadoTicket: "PRO",
		RutResponsable:  "",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ==================== Close (Cerrar ticket) ====================

func createTicketInTER(t *testing.T, svc *services.TicketService) domain.Ticket {
	t.Helper()
	ticket := createTicketInPRO(t, svc)
	ctx := context.Background()

	// PRO → TER
	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "TER",
		RutResponsable:  "12345678-K",
	})
	if err != nil {
		t.Fatalf("change estado TER: %v", err)
	}

	updated, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get ticket: %v", err)
	}
	return updated.Ticket
}

func TestCloseTicketOK(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInTER(t, svc)

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      ticket.ID,
		IDSolicitante: 1,
		Nota:          5,
		Comentarios:   "Excelente trabajo",
		Observacion:   "Se reparó equipo correctamente",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verificar estado CER
	var row models.Ticket
	if err := db.Where("id = ?", ticket.ID).Take(&row).Error; err != nil {
		t.Fatalf("ticket not found: %v", err)
	}
	if row.CodEstadoTicket != "CER" {
		t.Errorf("estado = %q, want CER", row.CodEstadoTicket)
	}

	// Verificar valorización
	var val models.Valorizacion
	if err := db.Where("id_ticket = ?", ticket.ID).Take(&val).Error; err != nil {
		t.Fatalf("valorizacion not found: %v", err)
	}
	if val.Nota != 5 {
		t.Errorf("nota = %d, want 5", val.Nota)
	}
	if val.IDTecnico != 1 {
		t.Errorf("id_tecnico = %d, want 1", val.IDTecnico)
	}
	if val.Comentarios != "Excelente trabajo" {
		t.Errorf("comentarios = %q, want 'Excelente trabajo'", val.Comentarios)
	}

	// Verificar bitácora creada con la observación
	var bit models.BitacoraTicket
	if err := db.Where("id_ticket = ? AND comentario = ?", ticket.ID, "Se reparó equipo correctamente").Take(&bit).Error; err != nil {
		t.Fatalf("bitacora not found: %v", err)
	}
	if bit.RutAutor != "12345678-K" {
		t.Errorf("rut_autor = %q, want 12345678-K", bit.RutAutor)
	}

	// Verificar trazabilidad CER
	var traz models.TrazabilidadTicket
	if err := db.Where("id_ticket = ? AND cod_estado_ticket = ?", ticket.ID, "CER").Take(&traz).Error; err != nil {
		t.Fatalf("trazabilidad CER not found: %v", err)
	}
}

func TestCloseTicketNotInTER(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      ticket.ID,
		IDSolicitante: 1,
		Nota:          5,
		Observacion:   "obs",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCloseTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      999,
		IDSolicitante: 1,
		Nota:          5,
		Observacion:   "obs",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestCloseTicketInvalidNota(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      1,
		IDSolicitante: 1,
		Nota:          6,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCloseTicketMissingSolicitante(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      1,
		IDSolicitante: 0,
		Nota:          5,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCloseTicketWithoutComentarios(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInTER(t, svc)

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      ticket.ID,
		IDSolicitante: 1,
		Nota:          3,
		Observacion:   "Cierre sin comentarios",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var val models.Valorizacion
	if err := db.Where("id_ticket = ?", ticket.ID).Take(&val).Error; err != nil {
		t.Fatalf("valorizacion not found: %v", err)
	}
	if val.Nota != 3 {
		t.Errorf("nota = %d, want 3", val.Nota)
	}
	if val.Comentarios != "" {
		t.Errorf("comentarios = %q, want empty", val.Comentarios)
	}
}

func TestCloseTicketMissingObservacion(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	err := svc.Close(ctx, services.CloseTicketCommand{
		IDTicket:      1,
		IDSolicitante: 1,
		Nota:          5,
		Observacion:   "   ",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ==================== Pausas ====================

func createTicketInPRO(t *testing.T, svc *services.TicketService) domain.Ticket {
	t.Helper()
	ticket := createTicketForAssign(t, svc)
	ctx := context.Background()

	_, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          ticket.ID,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err != nil {
		t.Fatalf("assign: %v", err)
	}

	// ASI → VITEC
	err = svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "VITEC",
		RutResponsable:  "12345678-K",
	})
	if err != nil {
		t.Fatalf("change estado VITEC: %v", err)
	}

	// VITEC → PRO
	err = svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "PRO",
		RutResponsable:  "12345678-K",
	})
	if err != nil {
		t.Fatalf("change estado PRO: %v", err)
	}

	updated, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get ticket: %v", err)
	}
	return updated.Ticket
}

func TestCreatePausaRequiereAutorizacion(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	pausa, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  1, // requiere_autorizacion = true
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pausa.EstadoPausa != "PENDIENTE" {
		t.Errorf("estado_pausa = %q, want PENDIENTE", pausa.EstadoPausa)
	}

	// ticket debe estar en PAU
	var row models.Ticket
	if err := db.Where("id = ?", ticket.ID).Take(&row).Error; err != nil {
		t.Fatalf("ticket not found: %v", err)
	}
	if row.CodEstadoTicket != "PAU" {
		t.Errorf("estado ticket = %q, want PAU", row.CodEstadoTicket)
	}
}

func TestCreatePausaAutoAprobada(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	pausa, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  2, // requiere_autorizacion = false
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pausa.EstadoPausa != "APROBADA" {
		t.Errorf("estado_pausa = %q, want APROBADA", pausa.EstadoPausa)
	}
}

func TestCreatePausaTicketNotInPRO(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc) // estado CRE

	_, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  1,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCreatePausaTecnicoNotAssigned(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	_, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 999,
		IDMotivoPausa:  1,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestResolverPausaAprobar(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	pausa, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  1,
	})
	if err != nil {
		t.Fatalf("create pausa: %v", err)
	}

	err = svc.ResolverPausa(ctx, services.ResolverPausaCommand{
		IDPausa:             pausa.ID,
		EstadoPausa:         "APROBADA",
		IDTecnicoAutorizado: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var row models.TicketPausa
	if err := db.Where("id = ?", pausa.ID).Take(&row).Error; err != nil {
		t.Fatalf("pausa not found: %v", err)
	}
	if row.EstadoPausa != "APROBADA" {
		t.Errorf("estado_pausa = %q, want APROBADA", row.EstadoPausa)
	}
	if row.IDTecnicoAutorizado == nil || *row.IDTecnicoAutorizado != 1 {
		t.Error("id_tecnico_autorizado should be 1")
	}
}

func TestResolverPausaRechazar(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	pausa, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  1,
	})
	if err != nil {
		t.Fatalf("create pausa: %v", err)
	}

	err = svc.ResolverPausa(ctx, services.ResolverPausaCommand{
		IDPausa:             pausa.ID,
		EstadoPausa:         "RECHAZADA",
		IDTecnicoAutorizado: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var row models.TicketPausa
	if err := db.Where("id = ?", pausa.ID).Take(&row).Error; err != nil {
		t.Fatalf("pausa not found: %v", err)
	}
	if row.EstadoPausa != "RECHAZADA" {
		t.Errorf("estado_pausa = %q, want RECHAZADA", row.EstadoPausa)
	}

	var ticketRow models.Ticket
	if err := db.Where("id = ?", ticket.ID).Take(&ticketRow).Error; err != nil {
		t.Fatalf("ticket not found: %v", err)
	}
	if ticketRow.CodEstadoTicket != "PRO" {
		t.Errorf("estado ticket = %q, want PRO", ticketRow.CodEstadoTicket)
	}
}

func TestResolverPausaYaResuelta(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	pausa, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  2, // auto-aprobada
	})
	if err != nil {
		t.Fatalf("create pausa: %v", err)
	}

	err = svc.ResolverPausa(ctx, services.ResolverPausaCommand{
		IDPausa:             pausa.ID,
		EstadoPausa:         "APROBADA",
		IDTecnicoAutorizado: 1,
	})
	if err == nil {
		t.Fatal("expected error, pausa already resolved")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestReanudarTicketOK(t *testing.T) {
	t.Parallel()
	svc, db := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	_, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  2, // auto-aprobada
	})
	if err != nil {
		t.Fatalf("create pausa: %v", err)
	}

	err = svc.ReanudarTicket(ctx, services.ReanudarTicketCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var row models.Ticket
	if err := db.Where("id = ?", ticket.ID).Take(&row).Error; err != nil {
		t.Fatalf("ticket not found: %v", err)
	}
	if row.CodEstadoTicket != "PRO" {
		t.Errorf("estado = %q, want PRO", row.CodEstadoTicket)
	}
}

func TestReanudarTicketNotPaused(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	err := svc.ReanudarTicket(ctx, services.ReanudarTicketCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
	})
	if err == nil {
		t.Fatal("expected error, ticket not paused")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── ListPausas ───

func TestListPausasAll(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	// crear pausa auto-aprobada
	_, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  2,
	})
	if err != nil {
		t.Fatalf("create pausa 1: %v", err)
	}

	// reanudar para poder crear otra pausa
	err = svc.ReanudarTicket(ctx, services.ReanudarTicketCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
	})
	if err != nil {
		t.Fatalf("reanudar: %v", err)
	}

	// crear pausa pendiente
	_, err = svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  1,
	})
	if err != nil {
		t.Fatalf("create pausa 2: %v", err)
	}

	// listar todas
	result, err := svc.ListPausas(ctx, services.ListPausasQuery{
		IDTicket: ticket.ID,
	})
	if err != nil {
		t.Fatalf("list pausas: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("total = %d, want 2", result.Total)
	}
	if len(result.Items) != 2 {
		t.Errorf("items = %d, want 2", len(result.Items))
	}
}

func TestListPausasFilterByEstado(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	// crear pausa auto-aprobada
	_, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  2,
	})
	if err != nil {
		t.Fatalf("create pausa: %v", err)
	}

	// reanudar
	err = svc.ReanudarTicket(ctx, services.ReanudarTicketCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
	})
	if err != nil {
		t.Fatalf("reanudar: %v", err)
	}

	// crear pausa pendiente
	_, err = svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  1,
	})
	if err != nil {
		t.Fatalf("create pausa 2: %v", err)
	}

	// filtrar solo PENDIENTE
	result, err := svc.ListPausas(ctx, services.ListPausasQuery{
		IDTicket: ticket.ID,
		Estado:   "PENDIENTE",
	})
	if err != nil {
		t.Fatalf("list pausas pendientes: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total = %d, want 1", result.Total)
	}
	if len(result.Items) != 1 {
		t.Errorf("items = %d, want 1", len(result.Items))
	}
	if result.Items[0].EstadoPausa != "PENDIENTE" {
		t.Errorf("estado = %s, want PENDIENTE", result.Items[0].EstadoPausa)
	}

	// filtrar solo APROBADA
	result, err = svc.ListPausas(ctx, services.ListPausasQuery{
		IDTicket: ticket.ID,
		Estado:   "APROBADA",
	})
	if err != nil {
		t.Fatalf("list pausas aprobadas: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total = %d, want 1", result.Total)
	}
}

func TestListPausasPagination(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	// crear pausa
	_, err := svc.CreatePausa(ctx, services.CreatePausaCommand{
		IDTicket:       ticket.ID,
		IDTecnicoPausa: 1,
		IDMotivoPausa:  2,
	})
	if err != nil {
		t.Fatalf("create pausa: %v", err)
	}

	// listar con offset que salta todo
	result, err := svc.ListPausas(ctx, services.ListPausasQuery{
		IDTicket: ticket.ID,
		Limit:    10,
		Offset:   10,
	})
	if err != nil {
		t.Fatalf("list pausas: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total = %d, want 1", result.Total)
	}
	if len(result.Items) != 0 {
		t.Errorf("items = %d, want 0 (offset beyond data)", len(result.Items))
	}
}

func TestListPausasTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.ListPausas(ctx, services.ListPausasQuery{
		IDTicket: 999,
	})
	if err == nil {
		t.Fatal("expected error, ticket not found")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

func TestListPausasInvalidEstado(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.ListPausas(ctx, services.ListPausasQuery{
		IDTicket: 1,
		Estado:   "INVALIDO",
	})
	if err == nil {
		t.Fatal("expected error, invalid estado")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── Traspasos ───

func TestCreateTraspasoOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	traspaso, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "No tengo experiencia en redes",
	})
	if err != nil {
		t.Fatalf("create traspaso: %v", err)
	}

	if traspaso.EstadoTraspaso != "PENDIENTE" {
		t.Errorf("estado = %s, want PENDIENTE", traspaso.EstadoTraspaso)
	}
	if traspaso.IDTecnicoOrigen != 1 {
		t.Errorf("tecnico_origen = %d, want 1", traspaso.IDTecnicoOrigen)
	}
	if traspaso.IDTecnicoDestino != 2 {
		t.Errorf("tecnico_destino = %d, want 2", traspaso.IDTecnicoDestino)
	}

	// verificar que el ticket cambió a STR
	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get ticket: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "STR" {
		t.Errorf("ticket estado = %s, want STR", detalle.Ticket.CodEstadoTicket)
	}

	// verificar trazabilidad STR
	found := false
	for _, traz := range detalle.Trazabilidad {
		if traz.CodEstadoTicket == "STR" {
			found = true
		}
	}
	if !found {
		t.Error("trazabilidad STR not found")
	}
}

func TestCreateTraspasoTicketNotInPRO(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc) // está en CRE

	_, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "motivo",
	})
	if err == nil {
		t.Fatal("expected error, ticket not in PRO")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCreateTraspasoTecnicoNotAssigned(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	_, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  999, // no es el asignado
		IDTecnicoDestino: 2,
		Motivo:           "motivo",
	})
	if err == nil {
		t.Fatal("expected error, tecnico not assigned")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCreateTraspasoSameTecnico(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         1,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 1,
		Motivo:           "motivo",
	})
	if err == nil {
		t.Fatal("expected error, same tecnico")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestCreateTraspasoDuplicatePendiente(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	// crear primer traspaso
	_, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "motivo",
	})
	if err != nil {
		t.Fatalf("create traspaso 1: %v", err)
	}

	// intentar crear otro (ticket ya está en STR, no en PRO)
	_, err = svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "otro motivo",
	})
	if err == nil {
		t.Fatal("expected error, duplicate pending traspaso")
	}
}

func TestResolverTraspasoAceptado(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	traspaso, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "Necesito especialista en redes",
	})
	if err != nil {
		t.Fatalf("create traspaso: %v", err)
	}

	err = svc.ResolverTraspaso(ctx, services.ResolverTraspasoCommand{
		IDTraspaso:           traspaso.ID,
		EstadoTraspaso:       "ACEPTADO",
		ComentarioResolucion: "Sin problema",
	})
	if err != nil {
		t.Fatalf("resolver traspaso: %v", err)
	}

	// verificar ticket: técnico cambiado + estado PRO
	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get ticket: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "PRO" {
		t.Errorf("ticket estado = %s, want PRO", detalle.Ticket.CodEstadoTicket)
	}
	if detalle.Ticket.IDTecnicoAsignado == nil || *detalle.Ticket.IDTecnicoAsignado != 2 {
		t.Errorf("tecnico_asignado = %v, want 2", detalle.Ticket.IDTecnicoAsignado)
	}

	// verificar trazabilidad: debe tener TRA y PRO
	foundTRA := false
	foundPRO := false
	for _, traz := range detalle.Trazabilidad {
		if traz.CodEstadoTicket == "TRA" {
			foundTRA = true
		}
		if traz.CodEstadoTicket == "PRO" {
			foundPRO = true
		}
	}
	if !foundTRA {
		t.Error("trazabilidad TRA not found")
	}
	if !foundPRO {
		t.Error("trazabilidad PRO not found after TRA")
	}
}

func TestResolverTraspasoRechazado(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	traspaso, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "motivo",
	})
	if err != nil {
		t.Fatalf("create traspaso: %v", err)
	}

	err = svc.ResolverTraspaso(ctx, services.ResolverTraspasoCommand{
		IDTraspaso:           traspaso.ID,
		EstadoTraspaso:       "RECHAZADO",
		ComentarioResolucion: "No puedo tomar el ticket",
	})
	if err != nil {
		t.Fatalf("resolver traspaso: %v", err)
	}

	// verificar ticket: técnico sigue siendo el mismo + estado PRO
	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get ticket: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "PRO" {
		t.Errorf("ticket estado = %s, want PRO", detalle.Ticket.CodEstadoTicket)
	}
	if detalle.Ticket.IDTecnicoAsignado == nil || *detalle.Ticket.IDTecnicoAsignado != 1 {
		t.Errorf("tecnico_asignado = %v, want 1 (unchanged)", detalle.Ticket.IDTecnicoAsignado)
	}
}

func TestResolverTraspasoYaResuelto(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	traspaso, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "motivo",
	})
	if err != nil {
		t.Fatalf("create traspaso: %v", err)
	}

	// resolver
	err = svc.ResolverTraspaso(ctx, services.ResolverTraspasoCommand{
		IDTraspaso:     traspaso.ID,
		EstadoTraspaso: "RECHAZADO",
	})
	if err != nil {
		t.Fatalf("resolver 1: %v", err)
	}

	// intentar resolver de nuevo
	err = svc.ResolverTraspaso(ctx, services.ResolverTraspasoCommand{
		IDTraspaso:     traspaso.ID,
		EstadoTraspaso: "ACEPTADO",
	})
	if err == nil {
		t.Fatal("expected error, traspaso already resolved")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

func TestListTraspasosWithFilter(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	// crear traspaso y rechazarlo
	traspaso, err := svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "primer intento",
	})
	if err != nil {
		t.Fatalf("create traspaso: %v", err)
	}

	err = svc.ResolverTraspaso(ctx, services.ResolverTraspasoCommand{
		IDTraspaso:     traspaso.ID,
		EstadoTraspaso: "RECHAZADO",
	})
	if err != nil {
		t.Fatalf("resolver: %v", err)
	}

	// crear otro traspaso (ticket volvió a PRO)
	_, err = svc.CreateTraspaso(ctx, services.CreateTraspasoCommand{
		IDTicket:         ticket.ID,
		IDTecnicoOrigen:  1,
		IDTecnicoDestino: 2,
		Motivo:           "segundo intento",
	})
	if err != nil {
		t.Fatalf("create traspaso 2: %v", err)
	}

	// listar todos
	result, err := svc.ListTraspasos(ctx, services.ListTraspasosQuery{
		IDTicket: ticket.ID,
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("total = %d, want 2", result.Total)
	}

	// filtrar PENDIENTE
	result, err = svc.ListTraspasos(ctx, services.ListTraspasosQuery{
		IDTicket: ticket.ID,
		Estado:   "PENDIENTE",
	})
	if err != nil {
		t.Fatalf("list pendiente: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total pendiente = %d, want 1", result.Total)
	}

	// filtrar RECHAZADO
	result, err = svc.ListTraspasos(ctx, services.ListTraspasosQuery{
		IDTicket: ticket.ID,
		Estado:   "RECHAZADO",
	})
	if err != nil {
		t.Fatalf("list rechazado: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("total rechazado = %d, want 1", result.Total)
	}
}

// ─── Cancelar vía ChangeEstado ───

func TestChangeEstadoCancelOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "CAN",
		RutResponsable:  "12345678-9",
	})
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}

	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "CAN" {
		t.Errorf("estado = %s, want CAN", detalle.Ticket.CodEstadoTicket)
	}

	found := false
	for _, traz := range detalle.Trazabilidad {
		if traz.CodEstadoTicket == "CAN" {
			found = true
		}
	}
	if !found {
		t.Error("trazabilidad CAN not found")
	}
}

func TestChangeEstadoCancelFromPRO(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc)

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "CAN",
		RutResponsable:  "12345678-9",
	})
	if err != nil {
		t.Fatalf("cancel from PRO: %v", err)
	}

	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "CAN" {
		t.Errorf("estado = %s, want CAN", detalle.Ticket.CodEstadoTicket)
	}
}

func TestChangeEstadoCancelAlreadyCancelled(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	_ = svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "CAN",
		RutResponsable:  "12345678-9",
	})

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "CAN",
		RutResponsable:  "12345678-9",
	})
	if err == nil {
		t.Fatal("expected error, already cancelled")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── Visto Técnico vía ChangeEstado ───

func TestChangeEstadoVistoOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	// asignar ticket
	assigned, err := svc.Assign(ctx, services.AssignTicketCommand{
		IDTicket:          ticket.ID,
		IDTecnicoAsignado: 1,
		IDCatalogoFalla:   1,
		IDNivelPrioridad:  1,
	})
	if err != nil {
		t.Fatalf("assign: %v", err)
	}

	err = svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        assigned.ID,
		CodEstadoTicket: "VITEC",
		RutResponsable:  "11111111-1",
	})
	if err != nil {
		t.Fatalf("visto: %v", err)
	}

	detalle, err := svc.GetByNroTicket(ctx, assigned.NroTicket)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "VITEC" {
		t.Errorf("estado = %s, want VITEC", detalle.Ticket.CodEstadoTicket)
	}
}

func TestChangeEstadoVistoNotASI(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc) // estado CRE

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "VITEC",
		RutResponsable:  "11111111-1",
	})
	if err == nil {
		t.Fatal("expected error, not in ASI")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── ChangeEstado TER ───

func TestChangeEstadoTerminarOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketInPRO(t, svc) // PRO

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "TER",
		RutResponsable:  "11111111-1",
	})
	if err != nil {
		t.Fatalf("terminar: %v", err)
	}

	detalle, err := svc.GetByNroTicket(ctx, ticket.NroTicket)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if detalle.Ticket.CodEstadoTicket != "TER" {
		t.Errorf("estado = %s, want TER", detalle.Ticket.CodEstadoTicket)
	}
}

func TestChangeEstadoTerminarNotPRO(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc) // CRE

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "TER",
		RutResponsable:  "12345678-9",
	})
	if err == nil {
		t.Fatal("expected error, not in PRO")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── ChangeEstado código no permitido ───

func TestChangeEstadoCodigoNoPermitido(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "ASI",
		RutResponsable:  "12345678-9",
	})
	if err == nil {
		t.Fatal("expected error, ASI not allowed via this endpoint")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── ChangeEstado no permite CER directo ───

func TestChangeEstadoCerrarBlocked(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	err := svc.ChangeEstado(ctx, services.ChangeEstadoCommand{
		IDTicket:        ticket.ID,
		CodEstadoTicket: "CER",
		RutResponsable:  "12345678-9",
	})
	if err == nil {
		t.Fatal("expected error, CER should use /cerrar endpoint")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── Update Ticket ───

func TestUpdateTicketOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	newDetalle := "NUEVO DETALLE DE FALLA"
	newUb := "PISO 3 SALA 5"
	critico := true

	updated, err := svc.UpdateTicket(ctx, services.UpdateTicketCommand{
		IDTicket:              ticket.ID,
		DetalleFallaReportada: &newDetalle,
		UbicacionObs:          &newUb,
		Critico:               &critico,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	if updated.DetalleFallaReportada != newDetalle {
		t.Errorf("detalle = %s, want %s", updated.DetalleFallaReportada, newDetalle)
	}
	if updated.UbicacionObs != "PISO 3 SALA 5" {
		t.Errorf("ubicacion = %s, want PISO 3 SALA 5", updated.UbicacionObs)
	}
	if !updated.Critico {
		t.Error("critico = false, want true")
	}
}

func TestUpdateTicketNoFields(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	_, err := svc.UpdateTicket(ctx, services.UpdateTicketCommand{
		IDTicket: ticket.ID,
	})
	if err == nil {
		t.Fatal("expected error, no fields")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", appErr.Status())
	}
}

// ─── Get by ID ───

func TestGetTicketByIDOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	got, err := svc.GetByID(ctx, ticket.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.NroTicket != ticket.NroTicket {
		t.Errorf("nro = %s, want %s", got.NroTicket, ticket.NroTicket)
	}
}

func TestGetTicketByIDNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 99999)
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

// ─── List Tickets ───

func TestListTicketsOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	// crear 3 tickets
	for i := 0; i < 3; i++ {
		createTicketForAssign(t, svc)
	}

	result, err := svc.ListTickets(ctx, services.ListTicketsQuery{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total < 3 {
		t.Errorf("total = %d, want >= 3", result.Total)
	}
}

func TestListTicketsFilterEstado(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	createTicketForAssign(t, svc)            // CRE
	ticket2 := createTicketInPRO(t, svc)      // PRO
	_ = ticket2

	result, err := svc.ListTickets(ctx, services.ListTicketsQuery{
		CodEstadoTicket: "PRO",
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total < 1 {
		t.Errorf("total PRO = %d, want >= 1", result.Total)
	}
	for _, item := range result.Items {
		if item.CodEstadoTicket != "PRO" {
			t.Errorf("got estado %s, want PRO", item.CodEstadoTicket)
		}
	}
}

// ─── List Bitácora ───

func TestListBitacoraOK(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	ticket := createTicketForAssign(t, svc)

	// crear bitácora
	_, err := svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   ticket.ID,
		RutAutor:   "12345678-9",
		Comentario: "Primer comentario",
	})
	if err != nil {
		t.Fatalf("create bitacora: %v", err)
	}

	_, err = svc.CreateBitacora(ctx, services.CreateBitacoraCommand{
		IDTicket:   ticket.ID,
		RutAutor:   "12345678-9",
		Comentario: "Segundo comentario",
	})
	if err != nil {
		t.Fatalf("create bitacora 2: %v", err)
	}

	items, err := svc.ListBitacora(ctx, ticket.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("len = %d, want 2", len(items))
	}
}

func TestListBitacoraTicketNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := setupTicketService(t)
	ctx := context.Background()

	_, err := svc.ListBitacora(ctx, 99999)
	if err == nil {
		t.Fatal("expected error, ticket not found")
	}
	var appErr *domain.Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected domain.Error")
	}
	if appErr.Status() != http.StatusNotFound {
		t.Errorf("status = %d, want 404", appErr.Status())
	}
}

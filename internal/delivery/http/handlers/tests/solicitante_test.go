package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"

	"soporte/internal/adapters/repository/models"
	"soporte/internal/adapters/repository/repository"
	"soporte/internal/application/services"
	"soporte/internal/delivery/http/handlers"
	"soporte/internal/delivery/http/middlewares"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newSolicitanteRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)

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
	h := handlers.NewSolicitanteHandler(svc)

	router := gin.New()
	router.Use(middlewares.RequestID())
	group := router.Group("/api/v1/solicitantes")
	group.GET("", h.List)
	group.GET("/rut/:rut", h.GetByRut)
	group.GET("/:id", h.Get)
	group.POST("", h.Create)
	group.PATCH("/:id", h.Update)
	return router, db
}

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

func createSolicitante(t *testing.T, router *gin.Engine, body map[string]any) int {
	t.Helper()
	b, _ := json.Marshal(body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/v1/solicitantes", bytes.NewReader(b)))
	if resp.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %s", resp.Code, resp.Body.String())
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	data := result["data"].(map[string]any)
	return int(data["id"].(float64))
}

// --- Create ---

func TestHandlerCreateSolicitante(t *testing.T) {
	t.Parallel()
	router, db := newSolicitanteRouter(t)
	idServicio := createServicio(t, db, "TI")

	body := map[string]any{
		"rut":             "12345678",
		"dv":              "9",
		"nombre_completo": "Juan Perez",
		"id_servicio":     *idServicio,
		"correo":          "Juan@Test.com",
		"anexo":           100,
	}
	b, _ := json.Marshal(body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/v1/solicitantes", bytes.NewReader(b)))

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]any)

	if data["id"] == nil {
		t.Fatal("expected created solicitante id")
	}

	location := resp.Header().Get("Location")
	if location == "" {
		t.Fatal("expected Location header")
	}
}

// --- Get ---

func TestHandlerGetSolicitante(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	id := createSolicitante(t, router, map[string]any{
		"rut": "12345678", "dv": "9", "nombre_completo": "Test",
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/solicitantes/%d", id), nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]any)
	if data["nombre_completo"] != "Test" {
		t.Fatalf("expected nombre_completo=Test, got %v", data["nombre_completo"])
	}
}

func TestHandlerGetSolicitanteNotFound(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes/9999", nil))

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}

func TestHandlerGetSolicitanteByRut(t *testing.T) {
	t.Parallel()
	router, db := newSolicitanteRouter(t)
	idServicio := createServicio(t, db, "TI")

	id := createSolicitante(t, router, map[string]any{
		"rut": "12345678", "dv": "9", "nombre_completo": "Test", "id_servicio": *idServicio,
	})
	if id == 0 {
		t.Fatal("expected created id")
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes/rut/12.345.678-9", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]any)
	if data["rut"] != "12345678" {
		t.Fatalf("expected rut=12345678, got %v", data["rut"])
	}
	servicio, ok := data["servicio"].(map[string]any)
	if !ok {
		t.Fatalf("expected servicio relation, got %v", data["servicio"])
	}
	if servicio["servicios"] != "TI" {
		t.Fatalf("expected servicio TI, got %v", servicio["servicios"])
	}
}

func TestHandlerGetSolicitanteByRutInvalid(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes/rut/12A", nil))

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestHandlerGetSolicitanteInvalidID(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes/abc", nil))

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

// --- List ---

func TestHandlerListSolicitantes(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	for i := range 3 {
		createSolicitante(t, router, map[string]any{
			"rut": fmt.Sprintf("1000000%d", i), "dv": "0", "nombre_completo": fmt.Sprintf("Persona %d", i),
		})
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 3 {
		t.Fatalf("expected 3 items, got %d", len(data))
	}
	meta := result["meta"].(map[string]any)
	if meta["total"].(float64) != 3 {
		t.Fatalf("expected total=3, got %v", meta["total"])
	}
}

func TestHandlerListSolicitantesWithSearch(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	createSolicitante(t, router, map[string]any{
		"rut": "11111111", "dv": "1", "nombre_completo": "Maria Lopez",
	})
	createSolicitante(t, router, map[string]any{
		"rut": "22222222", "dv": "2", "nombre_completo": "Pedro Garcia",
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes?q=maria", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(data))
	}
}

func TestHandlerListSolicitantesPagination(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	for i := range 5 {
		createSolicitante(t, router, map[string]any{
			"rut": fmt.Sprintf("3000000%d", i), "dv": "0", "nombre_completo": fmt.Sprintf("P%d", i),
		})
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes?limit=2&offset=1", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(data))
	}
	meta := result["meta"].(map[string]any)
	if meta["total"].(float64) != 5 {
		t.Fatalf("expected total=5, got %v", meta["total"])
	}
}

// --- Update ---

func TestHandlerUpdateSolicitante(t *testing.T) {
	t.Parallel()
	router, db := newSolicitanteRouter(t)
	idServicioInicial := createServicio(t, db, "TI")
	idServicioNuevo := createServicio(t, db, "RRHH")

	id := createSolicitante(t, router, map[string]any{
		"rut": "12345678", "dv": "9", "nombre_completo": "Original", "id_servicio": *idServicioInicial,
	})

	updateBody, _ := json.Marshal(map[string]any{
		"nombre_completo": "Actualizado",
		"id_servicio":     *idServicioNuevo,
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/solicitantes/%d", id), bytes.NewReader(updateBody)))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]any)
	if data["nombre_completo"] != "Actualizado" {
		t.Fatalf("expected Actualizado, got %v", data["nombre_completo"])
	}
	if data["id_servicio"] != float64(*idServicioNuevo) {
		t.Fatalf("expected id_servicio=%d, got %v", *idServicioNuevo, data["id_servicio"])
	}
	if data["rut"] != "12345678" {
		t.Fatalf("expected rut unchanged, got %v", data["rut"])
	}
}

func TestHandlerUpdateSolicitanteNotFound(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	body, _ := json.Marshal(map[string]any{"nombre_completo": "Test"})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPatch, "/api/v1/solicitantes/9999", bytes.NewReader(body)))

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}

func TestHandlerUpdateSolicitanteDeactivate(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	id := createSolicitante(t, router, map[string]any{
		"rut": "12345678", "dv": "9", "nombre_completo": "Test",
	})

	body, _ := json.Marshal(map[string]any{"estado": false})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/solicitantes/%d", id), bytes.NewReader(body)))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]any)
	if data["estado"] != false {
		t.Fatalf("expected estado=false, got %v", data["estado"])
	}
}

// --- Validation errors ---

func TestHandlerCreateSolicitanteValidation(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	tests := []struct {
		name string
		body map[string]any
	}{
		{"missing rut", map[string]any{"dv": "9", "nombre_completo": "Test"}},
		{"missing dv", map[string]any{"rut": "12345678", "nombre_completo": "Test"}},
		{"missing nombre", map[string]any{"rut": "12345678", "dv": "9"}},
		{"invalid email", map[string]any{"rut": "12345678", "dv": "9", "nombre_completo": "Test", "correo": "not-email"}},
		{"dv too long", map[string]any{"rut": "12345678", "dv": "99", "nombre_completo": "Test"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b, _ := json.Marshal(tc.body)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/v1/solicitantes", bytes.NewReader(b)))

			if resp.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", resp.Code, resp.Body.String())
			}

			var result map[string]any
			json.NewDecoder(resp.Body).Decode(&result)
			if result["error"] == nil {
				t.Fatal("expected error in response")
			}
		})
	}
}

func TestHandlerCreateSolicitanteMalformedJSON(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/v1/solicitantes", bytes.NewReader([]byte("{invalid"))))

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestHandlerCreateSolicitanteDuplicateRut(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	createSolicitante(t, router, map[string]any{
		"rut": "99999999", "dv": "K", "nombre_completo": "Original",
	})

	body, _ := json.Marshal(map[string]any{
		"rut": "99999999", "dv": "K", "nombre_completo": "Duplicate",
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/v1/solicitantes", bytes.NewReader(body)))

	if resp.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestHandlerResponseHasRequestID(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/solicitantes", nil))

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if result["request_id"] == nil || result["request_id"] == "" {
		t.Fatal("expected request_id in response")
	}
}

func TestHandlerCreateSolicitanteInvalidServicio(t *testing.T) {
	t.Parallel()
	router, _ := newSolicitanteRouter(t)

	body, _ := json.Marshal(map[string]any{
		"rut":             "12345678",
		"dv":              "9",
		"nombre_completo": "Test",
		"id_servicio":     9999,
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, "/api/v1/solicitantes", bytes.NewReader(body)))

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", resp.Code, resp.Body.String())
	}
}

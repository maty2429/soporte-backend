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

func newCatalogoRouter(t *testing.T) *gin.Engine {
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

	repo := repository.NewCatalogoRepository(db)
	svc := services.NewCatalogoService(repo)
	h := handlers.NewCatalogoHandler(svc)

	router := gin.New()
	router.Use(middlewares.RequestID())

	v1 := router.Group("/api/v1")

	tt := v1.Group("/tipos-ticket")
	tt.GET("", h.ListTiposTicket)
	tt.POST("", h.CreateTipoTicket)
	tt.PATCH("/:id", h.UpdateTipoTicket)

	np := v1.Group("/niveles-prioridad")
	np.GET("", h.ListNivelesPrioridad)
	np.POST("", h.CreateNivelPrioridad)
	np.PATCH("/:id", h.UpdateNivelPrioridad)

	tc := v1.Group("/tipos-tecnico")
	tc.GET("", h.ListTiposTecnico)
	tc.POST("", h.CreateTipoTecnico)
	tc.PATCH("/:id", h.UpdateTipoTecnico)

	ds := v1.Group("/departamentos-soporte")
	ds.GET("", h.ListDepartamentosSoporte)
	ds.POST("", h.CreateDepartamentoSoporte)
	ds.PATCH("/:id", h.UpdateDepartamentoSoporte)

	mp := v1.Group("/motivos-pausa")
	mp.GET("", h.ListMotivosPausa)
	mp.POST("", h.CreateMotivoPausa)
	mp.PATCH("/:id", h.UpdateMotivoPausa)

	tu := v1.Group("/tipos-turno")
	tu.GET("", h.ListTiposTurno)
	tu.POST("", h.CreateTipoTurno)
	tu.PATCH("/:id", h.UpdateTipoTurno)

	return router
}

func postJSON(t *testing.T, router *gin.Engine, url string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b)))
	return resp
}

func patchJSON(t *testing.T, router *gin.Engine, url string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(b)))
	return resp
}

func getData(t *testing.T, resp *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return result["data"].(map[string]any)
}

// --- List ---

func TestHandlerCatalogoList(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	postJSON(t, router, "/api/v1/niveles-prioridad", map[string]any{"descripcion": "Alta"})
	postJSON(t, router, "/api/v1/niveles-prioridad", map[string]any{"descripcion": "Media"})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/niveles-prioridad", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 2 {
		t.Fatalf("expected 2 items, got %d", len(data))
	}
}

func TestHandlerCatalogoListEmpty(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/tipos-tecnico", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]any)
	if len(data) != 0 {
		t.Fatalf("expected 0 items, got %d", len(data))
	}
}

// --- Create ---

func TestHandlerCatalogoCreate(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	resp := postJSON(t, router, "/api/v1/tipos-ticket", map[string]any{
		"cod_tipo_ticket": "INC",
		"descripcion":     "Incidencia",
	})

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", resp.Code, resp.Body.String())
	}

	location := resp.Header().Get("Location")
	if location == "" {
		t.Fatal("expected Location header")
	}

	data := getData(t, resp)
	if data["CodTipoTicket"] != "INC" {
		t.Fatalf("expected INC, got %v", data["CodTipoTicket"])
	}
}

func TestHandlerCatalogoCreateValidation(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	resp := postJSON(t, router, "/api/v1/niveles-prioridad", map[string]any{})

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestHandlerCatalogoCreateDuplicate(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	postJSON(t, router, "/api/v1/tipos-tecnico", map[string]any{"descripcion": "Especialista"})
	resp := postJSON(t, router, "/api/v1/tipos-tecnico", map[string]any{"descripcion": "Especialista"})

	if resp.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", resp.Code, resp.Body.String())
	}
}

// --- Update ---

func TestHandlerCatalogoUpdate(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	createResp := postJSON(t, router, "/api/v1/niveles-prioridad", map[string]any{"descripcion": "Alta"})
	created := getData(t, createResp)
	id := int(created["ID"].(float64))

	resp := patchJSON(t, router, fmt.Sprintf("/api/v1/niveles-prioridad/%d", id), map[string]any{"descripcion": "Muy Alta"})

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	data := getData(t, resp)
	if data["Descripcion"] != "Muy Alta" {
		t.Fatalf("expected Muy Alta, got %v", data["Descripcion"])
	}
}

func TestHandlerCatalogoUpdateNotFound(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	resp := patchJSON(t, router, "/api/v1/niveles-prioridad/9999", map[string]any{"descripcion": "Test"})

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}

func TestHandlerCatalogoResponseHasRequestID(t *testing.T) {
	t.Parallel()
	router := newCatalogoRouter(t)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/api/v1/niveles-prioridad", nil))

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if result["request_id"] == nil || result["request_id"] == "" {
		t.Fatal("expected request_id in response")
	}
}

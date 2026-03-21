package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"soporte/internal/config"
	"soporte/internal/delivery/http/handlers"
	"soporte/internal/delivery/http/middlewares"
)

func newHealthRouter(h handlers.HealthHandler) *gin.Engine {
	router := gin.New()
	router.Use(middlewares.RequestID())
	router.GET("/livez", h.Livez)
	router.GET("/readyz", h.Readyz)
	router.GET("/health", h.Get)
	return router
}

func TestHealthLivez(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	h := handlers.NewHealthHandler(config.Config{
		App: config.AppConfig{Name: "soporte", Version: "0.1.0"},
	}, nil, time.Now())

	resp := httptest.NewRecorder()
	newHealthRouter(h).ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/livez", nil))

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "alive" {
		t.Fatalf("expected status=alive, got %v", body["status"])
	}
}

func TestHealthReadyz(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		dbEnabled  bool
		wantCode   int
		wantStatus string
	}{
		{
			name:       "db deshabilitada → ready",
			dbEnabled:  false,
			wantCode:   http.StatusOK,
			wantStatus: "ready",
		},
		{
			name:       "db habilitada pero nil → not_ready",
			dbEnabled:  true,
			wantCode:   http.StatusServiceUnavailable,
			wantStatus: "not_ready",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.Config{
				App:      config.AppConfig{Name: "soporte", Version: "0.1.0"},
				Database: config.DatabaseConfig{Enabled: tc.dbEnabled},
			}
			h := handlers.NewHealthHandler(cfg, nil, time.Now())

			resp := httptest.NewRecorder()
			newHealthRouter(h).ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/readyz", nil))

			if resp.Code != tc.wantCode {
				t.Fatalf("expected %d, got %d", tc.wantCode, resp.Code)
			}

			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["status"] != tc.wantStatus {
				t.Fatalf("expected status=%s, got %v", tc.wantStatus, body["status"])
			}
		})
	}
}

func TestHealthGet(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		dbEnabled  bool
		wantCode   int
		wantStatus string
		wantDB     string
	}{
		{
			name:       "db deshabilitada → ok",
			dbEnabled:  false,
			wantCode:   http.StatusOK,
			wantStatus: "ok",
			wantDB:     "disabled",
		},
		{
			name:       "db habilitada pero nil → degraded",
			dbEnabled:  true,
			wantCode:   http.StatusServiceUnavailable,
			wantStatus: "degraded",
			wantDB:     "down",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.Config{
				App:      config.AppConfig{Name: "soporte", Version: "0.1.0"},
				Database: config.DatabaseConfig{Enabled: tc.dbEnabled},
			}
			h := handlers.NewHealthHandler(cfg, nil, time.Now())

			resp := httptest.NewRecorder()
			newHealthRouter(h).ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/health", nil))

			if resp.Code != tc.wantCode {
				t.Fatalf("expected %d, got %d", tc.wantCode, resp.Code)
			}

			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["status"] != tc.wantStatus {
				t.Fatalf("expected status=%s, got %v", tc.wantStatus, body["status"])
			}
			deps, _ := body["dependencies"].(map[string]any)
			if deps["database"] != tc.wantDB {
				t.Fatalf("expected database=%s, got %v", tc.wantDB, deps["database"])
			}
		})
	}
}

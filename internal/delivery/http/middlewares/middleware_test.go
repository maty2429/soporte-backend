package middlewares_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"soporte/internal/config"
	"soporte/internal/delivery/http/middlewares"
)

// ── ContentTypeJSON ───────────────────────────────────────────────────────────

func TestContentTypeJSON(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		wantCode    int
	}{
		{
			name:        "POST con application/json",
			method:      http.MethodPost,
			contentType: "application/json",
			body:        `{"x":1}`,
			wantCode:    http.StatusOK,
		},
		{
			name:        "POST con charset adicional",
			method:      http.MethodPost,
			contentType: "application/json; charset=utf-8",
			body:        `{"x":1}`,
			wantCode:    http.StatusOK,
		},
		{
			name:        "POST sin Content-Type",
			method:      http.MethodPost,
			contentType: "",
			body:        `{"x":1}`,
			wantCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:        "POST con text/plain",
			method:      http.MethodPost,
			contentType: "text/plain",
			body:        `{"x":1}`,
			wantCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:        "PUT sin Content-Type",
			method:      http.MethodPut,
			contentType: "",
			body:        `{"x":1}`,
			wantCode:    http.StatusUnsupportedMediaType,
		},
		{
			name:        "GET sin Content-Type pasa igual",
			method:      http.MethodGet,
			contentType: "",
			body:        "",
			wantCode:    http.StatusOK,
		},
		{
			name:        "DELETE sin Content-Type pasa igual",
			method:      http.MethodDelete,
			contentType: "",
			body:        "",
			wantCode:    http.StatusOK,
		},
		{
			name:        "POST sin body pasa igual",
			method:      http.MethodPost,
			contentType: "",
			body:        "",
			wantCode:    http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := gin.New()
			router.Use(middlewares.RequestID(), middlewares.ContentTypeJSON())
			router.Handle(tt.method, "/ok", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			var bodyReader *bytes.Reader
			if tt.body != "" {
				bodyReader = bytes.NewReader([]byte(tt.body))
			} else {
				bodyReader = bytes.NewReader(nil)
			}

			req := httptest.NewRequest(tt.method, "/ok", bodyReader)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.body != "" {
				req.ContentLength = int64(len(tt.body))
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.wantCode {
				t.Fatalf("esperaba %d, obtuvo %d — body: %s", tt.wantCode, resp.Code, resp.Body.String())
			}
		})
	}
}

// ── RequestID ─────────────────────────────────────────────────────────────────

func TestRequestID(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("genera ID cuando no se provee", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(middlewares.RequestID())
		router.GET("/ok", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		id := resp.Header().Get("X-Request-ID")
		if id == "" {
			t.Fatal("esperaba X-Request-ID en la respuesta, pero estaba vacío")
		}
	})

	t.Run("reutiliza el ID del header entrante", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(middlewares.RequestID())
		router.GET("/ok", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("X-Request-ID", "mi-id-personalizado")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		got := resp.Header().Get("X-Request-ID")
		if got != "mi-id-personalizado" {
			t.Fatalf("esperaba 'mi-id-personalizado', obtuvo %q", got)
		}
	})

	t.Run("ID disponible en contexto via GetRequestID", func(t *testing.T) {
		t.Parallel()

		var capturedID string

		router := gin.New()
		router.Use(middlewares.RequestID())
		router.GET("/ok", func(c *gin.Context) {
			capturedID = middlewares.GetRequestID(c)
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("X-Request-ID", "test-ctx-id")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if capturedID != "test-ctx-id" {
			t.Fatalf("esperaba 'test-ctx-id' en contexto, obtuvo %q", capturedID)
		}
	})

	t.Run("dos requests generan IDs distintos", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(middlewares.RequestID())
		router.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })

		resp1 := httptest.NewRecorder()
		router.ServeHTTP(resp1, httptest.NewRequest(http.MethodGet, "/ok", nil))

		resp2 := httptest.NewRecorder()
		router.ServeHTTP(resp2, httptest.NewRequest(http.MethodGet, "/ok", nil))

		id1 := resp1.Header().Get("X-Request-ID")
		id2 := resp2.Header().Get("X-Request-ID")

		if id1 == "" || id2 == "" {
			t.Fatal("ambos IDs deberían ser no vacíos")
		}
		if id1 == id2 {
			t.Fatalf("se esperaban IDs únicos, ambos fueron %q", id1)
		}
	})
}

// ── QueryTimeout ──────────────────────────────────────────────────────────────

func TestQueryTimeout(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("request completa antes del timeout", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.QueryTimeout(100*time.Millisecond))
		router.GET("/ok", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("esperaba 200, obtuvo %d", resp.Code)
		}
	})

	t.Run("request supera el timeout", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.QueryTimeout(5*time.Millisecond))
		router.GET("/slow", func(c *gin.Context) {
			// El handler espera a que el contexto sea cancelado
			// para no dejar goroutines colgadas en el test
			select {
			case <-time.After(5 * time.Second):
			case <-c.Request.Context().Done():
			}
		})

		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusGatewayTimeout {
			t.Fatalf("esperaba 504, obtuvo %d", resp.Code)
		}
	})

	t.Run("timeout desactivado con valor cero", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.QueryTimeout(0))
		router.GET("/ok", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("esperaba 200, obtuvo %d", resp.Code)
		}
	})
}

// ── Recovery ──────────────────────────────────────────────────────────────────

// TestRecovery no usa t.Parallel() en los subtests porque gin.CustomRecovery
// interactúa con el response writer durante el manejo del panic, lo que
// genera races con el detector cuando varios contextos gin corren en paralelo.
func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	t.Run("panic con string devuelve 500", func(t *testing.T) {
		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.Recovery(logger))
		router.GET("/panic", func(c *gin.Context) {
			panic("algo salió mal")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("esperaba 500, obtuvo %d", resp.Code)
		}

		body := resp.Body.String()
		if !bytes.Contains([]byte(body), []byte("internal_error")) {
			t.Fatalf("esperaba error code 'internal_error' en el body, obtuvo: %s", body)
		}
	})

	t.Run("panic con error devuelve 500", func(t *testing.T) {
		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.Recovery(logger))
		router.GET("/panic", func(c *gin.Context) {
			panic(http.ErrAbortHandler)
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("esperaba 500, obtuvo %d", resp.Code)
		}
	})

	t.Run("request siguiente funciona tras un panic", func(t *testing.T) {
		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.Recovery(logger))
		router.GET("/panic", func(c *gin.Context) { panic("boom") })
		router.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })

		// Primera request: panic
		panicResp := httptest.NewRecorder()
		router.ServeHTTP(panicResp, httptest.NewRequest(http.MethodGet, "/panic", nil))
		if panicResp.Code != http.StatusInternalServerError {
			t.Fatalf("esperaba 500 tras panic, obtuvo %d", panicResp.Code)
		}

		// Segunda request: normal — el servidor no debe haber crasheado
		okResp := httptest.NewRecorder()
		router.ServeHTTP(okResp, httptest.NewRequest(http.MethodGet, "/ok", nil))
		if okResp.Code != http.StatusOK {
			t.Fatalf("esperaba 200 tras recovery, obtuvo %d", okResp.Code)
		}
	})

	t.Run("respuesta incluye request_id", func(t *testing.T) {
		router := gin.New()
		router.Use(middlewares.RequestID(), middlewares.Recovery(logger))
		router.GET("/panic", func(c *gin.Context) { panic("boom") })

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		req.Header.Set("X-Request-ID", "test-recovery-id")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		body := resp.Body.String()
		if !bytes.Contains([]byte(body), []byte("test-recovery-id")) {
			t.Fatalf("esperaba request_id en la respuesta de error, obtuvo: %s", body)
		}
	})
}

// ── StructuredLogger ──────────────────────────────────────────────────────────

func TestStructuredLogger(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	newLoggerRouter := func(buf *bytes.Buffer) *gin.Engine {
		logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		r := gin.New()
		r.Use(middlewares.RequestID(), middlewares.StructuredLogger(logger))
		r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })
		r.GET("/fail", func(c *gin.Context) { c.Status(http.StatusInternalServerError) })
		r.GET("/warn", func(c *gin.Context) { c.Status(http.StatusBadRequest) })
		return r
	}

	t.Run("request exitosa loguea campos clave", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		router := newLoggerRouter(&buf)

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		router.ServeHTTP(httptest.NewRecorder(), req)

		logged := buf.String()
		for _, field := range []string{"method", "path", "status", "latency", "request_id"} {
			if !bytes.Contains([]byte(logged), []byte(field)) {
				t.Fatalf("esperaba campo %q en el log, log: %s", field, logged)
			}
		}
	})

	t.Run("query param seguro se incluye", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		router := newLoggerRouter(&buf)

		req := httptest.NewRequest(http.MethodGet, "/ok?limit=10", nil)
		router.ServeHTTP(httptest.NewRecorder(), req)

		if !bytes.Contains(buf.Bytes(), []byte("limit=10")) {
			t.Fatalf("esperaba 'limit=10' en el log, log: %s", buf.String())
		}
	})

	t.Run("query param sensible se redacta", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		router := newLoggerRouter(&buf)

		req := httptest.NewRequest(http.MethodGet, "/ok?token=secreto", nil)
		router.ServeHTTP(httptest.NewRecorder(), req)

		logged := buf.String()
		if bytes.Contains([]byte(logged), []byte("token=secreto")) {
			t.Fatalf("el param 'token' no debería aparecer en el log, log: %s", logged)
		}
		if !bytes.Contains([]byte(logged), []byte("_redacted=true")) {
			t.Fatalf("esperaba '_redacted=true' en el log, log: %s", logged)
		}
	})

	t.Run("status 5xx se loguea como error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		router := newLoggerRouter(&buf)

		req := httptest.NewRequest(http.MethodGet, "/fail", nil)
		router.ServeHTTP(httptest.NewRecorder(), req)

		if !bytes.Contains(buf.Bytes(), []byte(`"level":"ERROR"`)) {
			t.Fatalf("esperaba nivel ERROR para 5xx, log: %s", buf.String())
		}
	})

	t.Run("status 4xx se loguea como warn", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		router := newLoggerRouter(&buf)

		req := httptest.NewRequest(http.MethodGet, "/warn", nil)
		router.ServeHTTP(httptest.NewRecorder(), req)

		if !bytes.Contains(buf.Bytes(), []byte(`"level":"WARN"`)) {
			t.Fatalf("esperaba nivel WARN para 4xx, log: %s", buf.String())
		}
	})
}

// ── CORS ───────────────────────────────────────────────────────────────────────

func newCORSRouter(cfg config.SecurityConfig) *gin.Engine {
	r := gin.New()
	r.Use(middlewares.CORS(cfg))
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.OPTIONS("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestCORS(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("allow all devuelve wildcard", func(t *testing.T) {
		t.Parallel()

		router := newCORSRouter(config.SecurityConfig{CORSAllowAll: true})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "https://ejemplo.com")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "*" {
			t.Fatalf("esperaba '*', obtuvo %q", got)
		}
	})

	t.Run("origen en la lista se refleja", func(t *testing.T) {
		t.Parallel()

		router := newCORSRouter(config.SecurityConfig{
			CORSAllowedOrigins: []string{"https://app.ejemplo.com", "https://admin.ejemplo.com"},
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "https://app.ejemplo.com")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "https://app.ejemplo.com" {
			t.Fatalf("esperaba origen reflejado, obtuvo %q", got)
		}
		if vary := resp.Header().Get("Vary"); vary != "Origin" {
			t.Fatalf("esperaba Vary: Origin, obtuvo %q", vary)
		}
	})

	t.Run("origen no listado no recibe ACAO", func(t *testing.T) {
		t.Parallel()

		router := newCORSRouter(config.SecurityConfig{
			CORSAllowedOrigins: []string{"https://app.ejemplo.com"},
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "https://malicioso.com")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("esperaba ACAO vacío para origen no listado, obtuvo %q", got)
		}
	})

	t.Run("request sin origin no recibe ACAO", func(t *testing.T) {
		t.Parallel()

		router := newCORSRouter(config.SecurityConfig{
			CORSAllowedOrigins: []string{"https://app.ejemplo.com"},
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("esperaba ACAO vacío sin header Origin, obtuvo %q", got)
		}
	})

	t.Run("preflight OPTIONS devuelve 204 y headers CORS", func(t *testing.T) {
		t.Parallel()

		router := newCORSRouter(config.SecurityConfig{CORSAllowAll: true})

		req := httptest.NewRequest(http.MethodOptions, "/ok", nil)
		req.Header.Set("Origin", "https://ejemplo.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusNoContent {
			t.Fatalf("esperaba 204 en preflight, obtuvo %d", resp.Code)
		}
		if got := resp.Header().Get("Access-Control-Allow-Methods"); got == "" {
			t.Fatal("esperaba Access-Control-Allow-Methods en preflight")
		}
		if got := resp.Header().Get("Access-Control-Max-Age"); got == "" {
			t.Fatal("esperaba Access-Control-Max-Age en preflight")
		}
	})

	t.Run("segundo origen de la lista también se refleja", func(t *testing.T) {
		t.Parallel()

		router := newCORSRouter(config.SecurityConfig{
			CORSAllowedOrigins: []string{"https://app.ejemplo.com", "https://admin.ejemplo.com"},
		})

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set("Origin", "https://admin.ejemplo.com")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "https://admin.ejemplo.com" {
			t.Fatalf("esperaba segundo origen reflejado, obtuvo %q", got)
		}
	})
}

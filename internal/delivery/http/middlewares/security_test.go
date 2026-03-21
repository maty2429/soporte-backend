package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"soporte/internal/config"
	"soporte/internal/delivery/http/middlewares"
)

func TestSecurityHeaders(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.SecurityHeaders())
	router.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if got := resp.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("expected nosniff header, got %q", got)
	}
}

func TestRequestSizeLimit(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.RequestID(), middlewares.RequestSizeLimit(4))
	router.POST("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/ok", bytes.NewBufferString("12345"))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", resp.Code)
	}
}

func TestRateLimit(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middlewares.RequestID(), middlewares.RateLimit(config.SecurityConfig{
		RateLimitRequests: 1,
		RateLimitWindow:   time.Minute,
	}))
	router.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Primer request: permitido, headers presentes
	firstResp := httptest.NewRecorder()
	router.ServeHTTP(firstResp, httptest.NewRequest(http.MethodGet, "/ok", nil))

	if firstResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", firstResp.Code)
	}
	if firstResp.Header().Get("RateLimit-Limit") != "1" {
		t.Fatalf("expected RateLimit-Limit=1, got %q", firstResp.Header().Get("RateLimit-Limit"))
	}
	if firstResp.Header().Get("RateLimit-Remaining") != "0" {
		t.Fatalf("expected RateLimit-Remaining=0, got %q", firstResp.Header().Get("RateLimit-Remaining"))
	}

	// Segundo request: bloqueado, Retry-After presente
	secondResp := httptest.NewRecorder()
	router.ServeHTTP(secondResp, httptest.NewRequest(http.MethodGet, "/ok", nil))

	if secondResp.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", secondResp.Code)
	}
	if secondResp.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on 429")
	}
	if secondResp.Header().Get("RateLimit-Remaining") != "0" {
		t.Fatalf("expected RateLimit-Remaining=0 on 429, got %q", secondResp.Header().Get("RateLimit-Remaining"))
	}
}

package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"soporte/internal/adapters/repository/database"
	"soporte/internal/config"
	"soporte/internal/delivery/http/middlewares"
)

type HealthHandler struct {
	cfg       config.Config
	db        *gorm.DB
	startedAt time.Time
}

func NewHealthHandler(cfg config.Config, db *gorm.DB, startedAt time.Time) HealthHandler {
	return HealthHandler{
		cfg:       cfg,
		db:        db,
		startedAt: startedAt,
	}
}

func (h HealthHandler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var dbStatus string
	status := "ok"
	code := http.StatusOK

	if h.cfg.Database.Enabled && h.db == nil {
		dbStatus = "down"
		status = "degraded"
		code = http.StatusServiceUnavailable
	} else {
		err := database.Ping(ctx, h.db)
		switch {
		case err == nil:
			dbStatus = "up"
		case errors.Is(err, database.ErrDisabled):
			dbStatus = "disabled"
		default:
			dbStatus = "down"
			status = "degraded"
			code = http.StatusServiceUnavailable
		}
	}

	c.JSON(code, gin.H{
		"status":      status,
		"service":     h.cfg.App.Name,
		"version":     h.cfg.App.Version,
		"environment": h.cfg.App.Env,
		"time":        time.Now().UTC().Format(time.RFC3339),
		"uptime":      time.Since(h.startedAt).String(),
		"request_id":  middlewares.GetRequestID(c),
		"dependencies": gin.H{
			"database": dbStatus,
		},
	})
}

func (h HealthHandler) Livez(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "alive",
		"service":    h.cfg.App.Name,
		"version":    h.cfg.App.Version,
		"request_id": middlewares.GetRequestID(c),
	})
}

func (h HealthHandler) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if h.cfg.Database.Enabled && h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":     "not_ready",
			"service":    h.cfg.App.Name,
			"version":    h.cfg.App.Version,
			"request_id": middlewares.GetRequestID(c),
		})
		return
	}

	if err := database.Ping(ctx, h.db); err != nil && !errors.Is(err, database.ErrDisabled) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":     "not_ready",
			"service":    h.cfg.App.Name,
			"version":    h.cfg.App.Version,
			"request_id": middlewares.GetRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "ready",
		"service":    h.cfg.App.Name,
		"version":    h.cfg.App.Version,
		"request_id": middlewares.GetRequestID(c),
	})
}

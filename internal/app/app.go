package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	"soporte/internal/adapters/repository/database"
	"soporte/internal/config"
	httproutes "soporte/internal/delivery/http/routes"
)

type App struct {
	cfg    config.Config
	log    *slog.Logger
	server *http.Server
	db     *gorm.DB
}

func New(cfg config.Config, log *slog.Logger) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, cfg.Database, log)
	if err != nil {
		if cfg.Database.FailFast {
			return nil, err
		}

		fmt.Fprintf(os.Stderr, "\n\033[1;33m⚠  DATABASE UNAVAILABLE — starting in degraded mode\033[0m\n\033[33m   %v\033[0m\n\n", err)
		db = nil
	}

	router := httproutes.NewRouter(log, cfg, db, time.Now().UTC())

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           router,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadTimeout,
		MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
	}

	return &App{
		cfg:    cfg,
		log:    log,
		server: server,
		db:     db,
	}, nil
}

func (a *App) Start() error {
	dbStatus := "disabled"
	switch {
	case a.cfg.Database.Enabled && a.db != nil:
		dbStatus = "connected"
	case a.cfg.Database.Enabled && a.db == nil:
		dbStatus = "degraded"
	}

	a.log.Info("starting api server",
		"address", a.server.Addr,
		"environment", a.cfg.App.Env,
		"database_enabled", a.cfg.Database.Enabled,
		"database_status", dbStatus,
	)

	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("shutting down api server")

	serverErr := a.server.Shutdown(ctx)
	dbErr := database.Close(a.db)

	return errors.Join(serverErr, dbErr)
}

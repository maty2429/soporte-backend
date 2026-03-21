package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"soporte/internal/app"
	"soporte/internal/config"
)

func main() {
	envFile := "configs/.env.development"
	if config.IsProduction {
		envFile = "configs/.env.production"
	}

	if err := config.LoadEnvFiles(envFile); err != nil {
		slog.Error("load env files", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	// Si el tag de compilación es production, forzamos el entorno a production
	// e ignoramos cualquier valor de APP_ENV que venga de los archivos .env
	if config.IsProduction {
		cfg.App.Env = "production"
		// Desactivamos Swagger en producción independientemente de los archivos .env
		cfg.Docs.Enabled = false
	}

	logger := newLogger(cfg.App.Env)

	api, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("bootstrap application", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown application", "error", err)
		}
	}()

	if err := api.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("run server", "error", err)
		os.Exit(1)
	}
}

func newLogger(env string) *slog.Logger {
	var handler slog.Handler
	level := slog.LevelInfo

	// Si el tag de compilación es production o el env es production, usamos JSON y nivel Info
	if config.IsProduction || env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// En desarrollo/local usamos formato texto y nivel Debug
		if env == "development" || env == "local" || env == "" {
			level = slog.LevelDebug
		}
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}

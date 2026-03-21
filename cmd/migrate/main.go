package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"
	"time"

	"soporte/internal/adapters/repository/database"
	"soporte/internal/adapters/repository/migrations"
	"soporte/internal/config"
)

func main() {
	dir := flag.String("dir", migrations.DefaultDir, "directory containing migration files")
	steps := flag.Int("steps", 1, "number of migrations to rollback with down")
	flag.Parse()

	if err := config.LoadEnvFiles("configs/.env.development"); err != nil {
		slog.Error("load env files", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	if !cfg.Database.Enabled {
		slog.Error("database must be enabled to run migrations")
		os.Exit(1)
	}

	command := "status"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := database.Open(ctx, cfg.Database, slog.Default())
	if err != nil {
		slog.Error("open database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			slog.Error("close database", "error", err)
		}
	}()

	migrationsDir := migrations.MustDefaultDir(*dir)

	switch command {
	case "up":
		executed, err := migrations.Up(ctx, db, migrationsDir)
		if err != nil {
			slog.Error("run up migrations", "error", err)
			os.Exit(1)
		}

		fmt.Printf("applied %d migration(s)\n", executed)
	case "down":
		executed, err := migrations.Down(ctx, db, migrationsDir, *steps)
		if err != nil {
			slog.Error("run down migrations", "error", err)
			os.Exit(1)
		}

		fmt.Printf("rolled back %d migration(s)\n", executed)
	case "status":
		statuses, err := migrations.Statuses(ctx, db, migrationsDir)
		if err != nil {
			slog.Error("read migration status", "error", err)
			os.Exit(1)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(writer, "VERSION\tNAME\tSTATUS\tAPPLIED_AT")
		for _, item := range statuses {
			status := "pending"
			appliedAt := "-"
			if item.Applied {
				status = "applied"
				appliedAt = item.AppliedAt.UTC().Format(time.RFC3339)
			}

			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", item.Version, item.Name, status, appliedAt)
		}
		if err := writer.Flush(); err != nil {
			slog.Error("flush output", "error", err)
		}
	default:
		slog.Error("unknown migrate command", "command", command)
		os.Exit(1)
	}
}

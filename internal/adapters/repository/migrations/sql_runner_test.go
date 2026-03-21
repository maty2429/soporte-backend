package migrations_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"soporte/internal/adapters/repository/migrations"
)

func TestUpDownAndStatus(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeMigration(t, dir, "000001_create_items.up.sql", `
CREATE TABLE items (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);
`)
	writeMigration(t, dir, "000001_create_items.down.sql", `DROP TABLE IF EXISTS items;`)

	db := openSQLite(t)
	ctx := context.Background()

	executed, err := migrations.Up(ctx, db, dir)
	if err != nil {
		t.Fatalf("up migrations: %v", err)
	}

	if executed != 1 {
		t.Fatalf("expected 1 executed migration, got %d", executed)
	}

	statuses, err := migrations.Statuses(ctx, db, dir)
	if err != nil {
		t.Fatalf("statuses after up: %v", err)
	}

	if len(statuses) != 1 || !statuses[0].Applied {
		t.Fatalf("expected migration applied, got %+v", statuses)
	}

	executed, err = migrations.Down(ctx, db, dir, 1)
	if err != nil {
		t.Fatalf("down migrations: %v", err)
	}

	if executed != 1 {
		t.Fatalf("expected 1 rolled back migration, got %d", executed)
	}

	statuses, err = migrations.Statuses(ctx, db, dir)
	if err != nil {
		t.Fatalf("statuses after down: %v", err)
	}

	if len(statuses) != 1 || statuses[0].Applied {
		t.Fatalf("expected migration pending, got %+v", statuses)
	}
}

func TestLoadRequiresUpAndDown(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeMigration(t, dir, "000001_create_items.up.sql", `CREATE TABLE items (id INTEGER PRIMARY KEY);`)

	db := openSQLite(t)
	_, err := migrations.Up(context.Background(), db, dir)
	if err == nil {
		t.Fatal("expected error when down migration is missing")
	}
}

func writeMigration(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write migration %s: %v", name, err)
	}
}

func openSQLite(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	return db
}

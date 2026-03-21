package migrations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

const DefaultDir = "db/migrations"

type Migration struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

type Status struct {
	Version   string
	Name      string
	Applied   bool
	AppliedAt *time.Time
}

type appliedMigration struct {
	Version   string
	AppliedAt time.Time
}

func Up(ctx context.Context, db *gorm.DB, dir string) (int, error) {
	migrations, err := load(dir)
	if err != nil {
		return 0, err
	}

	if err := ensureTable(ctx, db); err != nil {
		return 0, err
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return 0, err
	}

	appliedSet := make(map[string]struct{}, len(applied))
	for _, item := range applied {
		appliedSet[item.Version] = struct{}{}
	}

	var executed int
	for _, migration := range migrations {
		if _, ok := appliedSet[migration.Version]; ok {
			continue
		}

		if err := applyUp(ctx, db, migration); err != nil {
			return executed, err
		}

		executed++
	}

	return executed, nil
}

func Down(ctx context.Context, db *gorm.DB, dir string, steps int) (int, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps must be greater than 0")
	}

	migrations, err := load(dir)
	if err != nil {
		return 0, err
	}

	if err := ensureTable(ctx, db); err != nil {
		return 0, err
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return 0, err
	}

	orderedApplied := lastApplied(migrations, applied)
	if len(orderedApplied) == 0 {
		return 0, nil
	}

	if steps > len(orderedApplied) {
		steps = len(orderedApplied)
	}

	var executed int
	for i := 0; i < steps; i++ {
		migration := orderedApplied[len(orderedApplied)-1-i]
		if err := applyDown(ctx, db, migration); err != nil {
			return executed, err
		}

		executed++
	}

	return executed, nil
}

func Statuses(ctx context.Context, db *gorm.DB, dir string) ([]Status, error) {
	migrations, err := load(dir)
	if err != nil {
		return nil, err
	}

	if err := ensureTable(ctx, db); err != nil {
		return nil, err
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]time.Time, len(applied))
	for _, item := range applied {
		appliedMap[item.Version] = item.AppliedAt
	}

	statuses := make([]Status, 0, len(migrations))
	for _, migration := range migrations {
		status := Status{
			Version: migration.Version,
			Name:    migration.Name,
		}

		if appliedAt, ok := appliedMap[migration.Version]; ok {
			status.Applied = true
			appliedCopy := appliedAt
			status.AppliedAt = &appliedCopy
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

func load(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	type pair struct {
		version  string
		name     string
		upPath   string
		downPath string
	}

	pairs := map[string]*pair{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		version, name, kind, ok := parseMigrationFilename(filename)
		if !ok {
			continue
		}

		key := version + ":" + name
		item, exists := pairs[key]
		if !exists {
			item = &pair{
				version: version,
				name:    name,
			}
			pairs[key] = item
		}

		fullPath := filepath.Join(dir, filename)
		switch kind {
		case "up":
			item.upPath = fullPath
		case "down":
			item.downPath = fullPath
		}
	}

	migrations := make([]Migration, 0, len(pairs))
	for _, item := range pairs {
		if item.upPath == "" || item.downPath == "" {
			return nil, fmt.Errorf("migration %s_%s must have both up and down files", item.version, item.name)
		}

		upSQL, err := os.ReadFile(item.upPath)
		if err != nil {
			return nil, fmt.Errorf("read up migration %s: %w", item.upPath, err)
		}

		downSQL, err := os.ReadFile(item.downPath)
		if err != nil {
			return nil, fmt.Errorf("read down migration %s: %w", item.downPath, err)
		}

		migrations = append(migrations, Migration{
			Version: item.version,
			Name:    item.name,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func parseMigrationFilename(filename string) (version, name, kind string, ok bool) {
	switch {
	case strings.HasSuffix(filename, ".up.sql"):
		kind = "up"
		filename = strings.TrimSuffix(filename, ".up.sql")
	case strings.HasSuffix(filename, ".down.sql"):
		kind = "down"
		filename = strings.TrimSuffix(filename, ".down.sql")
	default:
		return "", "", "", false
	}

	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", false
	}

	return parts[0], parts[1], kind, true
}

func ensureTable(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

func appliedVersions(ctx context.Context, db *gorm.DB) ([]appliedMigration, error) {
	type row struct {
		Version   string
		AppliedAt time.Time
	}

	var rows []row
	if err := db.WithContext(ctx).
		Table("schema_migrations").
		Select("version", "applied_at").
		Order("version ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list applied migrations: %w", err)
	}

	items := make([]appliedMigration, 0, len(rows))
	for _, r := range rows {
		items = append(items, appliedMigration(r))
	}

	return items, nil
}

func applyUp(ctx context.Context, db *gorm.DB, migration Migration) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(migration.UpSQL).Error; err != nil {
			return fmt.Errorf("apply migration %s_%s: %w", migration.Version, migration.Name, err)
		}

		if err := tx.Exec(
			`INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)`,
			migration.Version,
			migration.Name,
			time.Now().UTC(),
		).Error; err != nil {
			return fmt.Errorf("record migration %s_%s: %w", migration.Version, migration.Name, err)
		}

		return nil
	})
}

func applyDown(ctx context.Context, db *gorm.DB, migration Migration) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(migration.DownSQL).Error; err != nil {
			return fmt.Errorf("rollback migration %s_%s: %w", migration.Version, migration.Name, err)
		}

		if err := tx.Exec(
			`DELETE FROM schema_migrations WHERE version = ?`,
			migration.Version,
		).Error; err != nil {
			return fmt.Errorf("delete migration record %s_%s: %w", migration.Version, migration.Name, err)
		}

		return nil
	})
}

func lastApplied(migrations []Migration, applied []appliedMigration) []Migration {
	appliedSet := make(map[string]struct{}, len(applied))
	for _, item := range applied {
		appliedSet[item.Version] = struct{}{}
	}

	items := make([]Migration, 0, len(applied))
	for _, migration := range migrations {
		if _, ok := appliedSet[migration.Version]; ok {
			items = append(items, migration)
		}
	}

	return items
}

func MustDefaultDir(dir string) string {
	if strings.TrimSpace(dir) == "" {
		return DefaultDir
	}

	return dir
}

func Run(db *gorm.DB) error {
	if db == nil {
		return errors.New("database is required")
	}

	_, err := Up(context.Background(), db, DefaultDir)
	return err
}

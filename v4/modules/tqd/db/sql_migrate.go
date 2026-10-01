// Package db owns TQD SQL migration checks for v4 processes that open database.tqd.dsn.
//
// Canonical migration authority remains tqd-service/migrate (golang-migrate).
package db

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var migrationVersionPrefix = regexp.MustCompile(`^(\d+)_.*\.up\.sql$`)

type pendingMigration struct {
	version uint
	path    string
	name    string
}

// RunPendingMigrations applies newer *.up.sql files from tqd-service/migrate
// against the TQD database (schema_migrations golang-migrate layout).
func RunPendingMigrations(database *gorm.DB) error {
	if database == nil {
		return fmt.Errorf("tqd migrate database is nil")
	}

	dir, err := resolveMigrationDir()
	if err != nil {
		return err
	}

	if err := ensureGolangMigrateTable(database); err != nil {
		return err
	}

	current, dirty, err := readMigrationVersion(database)
	if err != nil {
		return err
	}
	if dirty {
		return fmt.Errorf("tqd schema_migrations is dirty at version %d; fix manually before start", current)
	}

	pending, err := listPendingMigrations(dir, current)
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		slog.Info(
			"TQD schema migrations up to date",
			"category", "lifecycle",
			"event_name", "database.schema.migrate",
			"schema_service", "tqd",
			"version", current,
			"migration_dir", dir,
		)
		return nil
	}

	for _, migration := range pending {
		slog.Info(
			"Applying TQD schema migration",
			"category", "lifecycle",
			"event_name", "database.schema.migrate.apply",
			"schema_service", "tqd",
			"version", migration.version,
			"file", migration.name,
		)
		if err := applyMigration(database, migration); err != nil {
			return err
		}
		current = migration.version
	}

	slog.Info(
		"TQD schema migrations applied",
		"category", "lifecycle",
		"event_name", "database.schema.migrate.done",
		"schema_service", "tqd",
		"version", current,
	)
	return nil
}

func resolveMigrationDir() (string, error) {
	candidates := []string{
		strings.TrimSpace(os.Getenv("TQD_MIGRATION_DIR")),
		filepath.Join("..", "tqd-service", "migrate"), // v4 cwd
		filepath.Join("tqd-service", "migrate"),        // repo root cwd
		"migrate",
		"/migrate",
		"/app/migrate",
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates,
			filepath.Join(filepath.Dir(exe), "migrate"),
			filepath.Join(filepath.Dir(exe), "..", "tqd-service", "migrate"),
		)
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			abs, err := filepath.Abs(candidate)
			if err != nil {
				return candidate, nil
			}
			return abs, nil
		}
	}
	return "", fmt.Errorf("tqd migrate directory not found (set TQD_MIGRATION_DIR or run from v4/ / repo root)")
}

func ensureGolangMigrateTable(database *gorm.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty BOOLEAN NOT NULL
)`
	if err := database.Exec(ddl).Error; err != nil {
		return fmt.Errorf("ensure tqd schema_migrations: %w", err)
	}
	return nil
}

func readMigrationVersion(database *gorm.DB) (uint, bool, error) {
	var count int64
	if err := database.Raw(`SELECT COUNT(1) FROM schema_migrations`).Scan(&count).Error; err != nil {
		return 0, false, fmt.Errorf("count tqd schema_migrations: %w", err)
	}
	if count == 0 {
		return 0, false, nil
	}

	var row struct {
		Version uint
		Dirty   bool
	}
	if err := database.Raw(`SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&row).Error; err != nil {
		return 0, false, fmt.Errorf("read tqd schema_migrations: %w", err)
	}
	return row.Version, row.Dirty, nil
}

func listPendingMigrations(dir string, current uint) ([]pendingMigration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read tqd migrate dir %s: %w", dir, err)
	}

	var pending []pendingMigration
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		match := migrationVersionPrefix.FindStringSubmatch(name)
		if match == nil {
			continue
		}
		version64, err := strconv.ParseUint(match[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse tqd migration version %s: %w", name, err)
		}
		version := uint(version64)
		if version <= current {
			continue
		}
		pending = append(pending, pendingMigration{
			version: version,
			path:    filepath.Join(dir, name),
			name:    name,
		})
	}
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].version < pending[j].version
	})
	return pending, nil
}

func applyMigration(database *gorm.DB, migration pendingMigration) error {
	body, err := os.ReadFile(migration.path)
	if err != nil {
		return fmt.Errorf("read %s: %w", migration.name, err)
	}

	if err := markMigrationDirty(database, migration.version); err != nil {
		return err
	}

	if err := database.Exec(string(body)).Error; err != nil {
		return fmt.Errorf("apply %s: %w", migration.name, err)
	}

	if err := markMigrationClean(database, migration.version); err != nil {
		return err
	}
	return nil
}

func markMigrationDirty(database *gorm.DB, version uint) error {
	var count int64
	if err := database.Raw(`SELECT COUNT(1) FROM schema_migrations`).Scan(&count).Error; err != nil {
		return fmt.Errorf("count schema_migrations before dirty mark: %w", err)
	}
	if count == 0 {
		if err := database.Exec(
			`INSERT INTO schema_migrations(version, dirty) VALUES (?, TRUE)`,
			version,
		).Error; err != nil {
			return fmt.Errorf("insert dirty tqd migration %d: %w", version, err)
		}
		return nil
	}
	if err := database.Exec(
		`UPDATE schema_migrations SET version = ?, dirty = TRUE`,
		version,
	).Error; err != nil {
		return fmt.Errorf("mark dirty tqd migration %d: %w", version, err)
	}
	return nil
}

func markMigrationClean(database *gorm.DB, version uint) error {
	if err := database.Exec(
		`UPDATE schema_migrations SET version = ?, dirty = FALSE`,
		version,
	).Error; err != nil {
		return fmt.Errorf("mark clean tqd migration %d: %w", version, err)
	}
	return nil
}

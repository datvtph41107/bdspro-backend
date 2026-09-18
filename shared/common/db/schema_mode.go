package _db

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SchemaMode selects the only schema mutation mechanism a process may use.
//
// SQL means the canonical service-local migrate/ runner owns schema changes;
// the serving process performs no DDL. AutoMigrate is an explicit bootstrap
// mode for local/development installations that need GORM to materialize the
// service's model-owned tables. Off is useful for tests and read-only tools.
type SchemaMode string

const (
	SchemaModeSQL         SchemaMode = "sql"
	SchemaModeAutoMigrate SchemaMode = "automigrate"
	SchemaModeOff         SchemaMode = "off"

	SchemaModeEnv = "QHPRO_DB_SCHEMA_MODE"
)

// SchemaPolicy is resolved once by the process root before it composes any
// repository. Source records which environment key selected the policy so an
// operator can prove which schema mechanism was active for a process.
type SchemaPolicy struct {
	Mode   SchemaMode
	Source string
}

// LoadSchemaPolicy resolves a service-specific override first, then the shared
// QHPRO_DB_SCHEMA_MODE key. The safe default is SQL: migrations must have been
// run by the service-local migrate/ owner before the serving process starts.
func LoadSchemaPolicy(service string) (SchemaPolicy, error) {
	serviceKey := serviceSchemaModeEnv(service)
	if raw, ok := os.LookupEnv(serviceKey); ok && strings.TrimSpace(raw) != "" {
		mode, err := ParseSchemaMode(raw)
		if err != nil {
			return SchemaPolicy{}, fmt.Errorf("parse %s: %w", serviceKey, err)
		}
		return SchemaPolicy{Mode: mode, Source: serviceKey}, nil
	}
	if raw, ok := os.LookupEnv(SchemaModeEnv); ok && strings.TrimSpace(raw) != "" {
		mode, err := ParseSchemaMode(raw)
		if err != nil {
			return SchemaPolicy{}, fmt.Errorf("parse %s: %w", SchemaModeEnv, err)
		}
		return SchemaPolicy{Mode: mode, Source: SchemaModeEnv}, nil
	}
	return SchemaPolicy{Mode: SchemaModeSQL, Source: "default"}, nil
}

func ParseSchemaMode(raw string) (SchemaMode, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "sql", "migrate", "migration":
		return SchemaModeSQL, nil
	case "automigrate", "auto-migrate", "gorm":
		return SchemaModeAutoMigrate, nil
	case "off", "disabled", "none":
		return SchemaModeOff, nil
	default:
		return "", fmt.Errorf(
			"unsupported schema mode %q (want sql, automigrate, or off)",
			raw,
		)
	}
}

// EnsureSchema is the process-root entry point for schema ownership: load
// policy, emit one lifecycle line, then apply. AutoMigrate runs with a Silent
// SQL logger so hundreds of information_schema probes do not flood the console;
// DDL failures still return to the process root.
//
// Prefer schema mode "sql" for day-to-day starts (zero DDL introspection). Use
// "automigrate" only for bootstrap of empty/local databases.
func EnsureSchema(
	database *gorm.DB,
	service string,
	autoMigrate func(*gorm.DB) error,
) error {
	policy, err := LoadSchemaPolicy(service)
	if err != nil {
		return fmt.Errorf("load %s schema policy: %w", service, err)
	}
	slog.Info(
		"Database schema policy",
		"category", "lifecycle",
		"event_name", "database.schema.policy",
		"schema_service", service,
		"mode", string(policy.Mode),
		"source", policy.Source,
	)
	migrateDB := database
	if policy.Mode == SchemaModeAutoMigrate && database != nil && database.Config != nil && database.Logger != nil {
		migrateDB = database.Session(&gorm.Session{
			Logger: database.Logger.LogMode(logger.Silent),
		})
	}
	if err := ApplySchemaPolicy(migrateDB, policy, autoMigrate); err != nil {
		return fmt.Errorf("apply %s schema policy: %w", service, err)
	}
	return nil
}

// ApplySchemaPolicy applies at most one schema mechanism. SQL is deliberately
// a no-op here: golang-migrate runs as the single deployment/init owner and
// serving binaries never attempt to replay its files. AutoMigrate invokes the
// service-owned model registry and propagates every error to the process root.
// Prefer EnsureSchema at process roots; keep this for tests and advanced callers.
func ApplySchemaPolicy(
	database *gorm.DB,
	policy SchemaPolicy,
	autoMigrate func(*gorm.DB) error,
) error {
	switch policy.Mode {
	case SchemaModeSQL, SchemaModeOff:
		return nil
	case SchemaModeAutoMigrate:
		if database == nil {
			return fmt.Errorf("automigrate database is nil")
		}
		if autoMigrate == nil {
			return fmt.Errorf("automigrate registry is nil")
		}
		if err := autoMigrate(database); err != nil {
			return fmt.Errorf("gorm automigrate: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported schema mode %q", policy.Mode)
	}
}

func serviceSchemaModeEnv(service string) string {
	normalized := strings.NewReplacer("-", "_", ".", "_").Replace(
		strings.ToUpper(strings.TrimSpace(service)),
	)
	if normalized == "" {
		return SchemaModeEnv
	}
	return "QHPRO_" + normalized + "_DB_SCHEMA_MODE"
}

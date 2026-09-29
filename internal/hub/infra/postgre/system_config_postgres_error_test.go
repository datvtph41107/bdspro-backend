package postgres

import (
	"context"
	"errors"
	"testing"

	"hub/internal/domain"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newSystemConfigDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		gormpostgres.New(gormpostgres.Config{
			DSN:                  "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			DryRun:                 true,
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}
	return db
}

func TestSystemConfigGetByKeyNormalizesRecordNotFound(t *testing.T) {
	db := newSystemConfigDryRunDB(t)
	if err := db.Callback().Query().Before("gorm:query").Register("test:not_found", func(tx *gorm.DB) {
		tx.AddError(gorm.ErrRecordNotFound)
	}); err != nil {
		t.Fatalf("register query callback: %v", err)
	}

	repo := &SystemConfigPostgres{DB: db}
	got, err := repo.GetByKey(context.Background(), "missing-key")
	if err != nil {
		t.Fatalf("GetByKey() error = %v, want nil", err)
	}
	if got != nil {
		t.Fatalf("GetByKey() entity = %#v, want nil", got)
	}
}

func TestSystemConfigBulkUpsertStopsOnLookupFailure(t *testing.T) {
	db := newSystemConfigDryRunDB(t)
	wantErr := errors.New("database unavailable")
	if err := db.Callback().Query().Before("gorm:query").Register("test:lookup_failure", func(tx *gorm.DB) {
		tx.AddError(wantErr)
	}); err != nil {
		t.Fatalf("register query callback: %v", err)
	}

	createCalls := 0
	if err := db.Callback().Create().Before("gorm:create").Register("test:count_create", func(*gorm.DB) {
		createCalls++
	}); err != nil {
		t.Fatalf("register create callback: %v", err)
	}

	repo := &SystemConfigPostgres{DB: db}
	got, created, updated, err := repo.BulkUpsert(context.Background(), []*domain.SystemConfigEntity{{Key: "config-key"}})
	if got != nil {
		t.Fatalf("BulkUpsert() results = %#v, want nil", got)
	}
	if created != 0 || updated != 0 {
		t.Fatalf("BulkUpsert() counts = (%d,%d), want (0,0)", created, updated)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("BulkUpsert() error = %v, want lookup error %v", err, wantErr)
	}
	if createCalls != 0 {
		t.Fatalf("BulkUpsert() create calls = %d, want 0", createCalls)
	}
}

func TestSystemConfigBulkUpsertCreatesOnlyAfterNormalizedAbsence(t *testing.T) {
	db := newSystemConfigDryRunDB(t)
	if err := db.Callback().Query().Before("gorm:query").Register("test:not_found", func(tx *gorm.DB) {
		tx.AddError(gorm.ErrRecordNotFound)
	}); err != nil {
		t.Fatalf("register query callback: %v", err)
	}

	createCalls := 0
	if err := db.Callback().Create().Before("gorm:create").Register("test:count_create", func(*gorm.DB) {
		createCalls++
	}); err != nil {
		t.Fatalf("register create callback: %v", err)
	}

	repo := &SystemConfigPostgres{DB: db}
	got, created, updated, err := repo.BulkUpsert(context.Background(), []*domain.SystemConfigEntity{{Key: "config-key"}})
	if err != nil {
		t.Fatalf("BulkUpsert() error = %v, want nil", err)
	}
	if len(got) != 1 {
		t.Fatalf("BulkUpsert() results len = %d, want 1", len(got))
	}
	if created != 1 || updated != 0 {
		t.Fatalf("BulkUpsert() counts = (%d,%d), want (1,0)", created, updated)
	}
	if createCalls != 1 {
		t.Fatalf("BulkUpsert() create calls = %d, want 1", createCalls)
	}
}

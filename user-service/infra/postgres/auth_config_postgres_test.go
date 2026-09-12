package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newAuthConfigDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  "host=localhost user=test dbname=test sslmode=disable",
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

func TestUpdateValuesAtomicallyPropagatesSecondWriteFailure(t *testing.T) {
	db := newAuthConfigDryRunDB(t)
	writeErr := errors.New("second update failed")
	calls := 0
	if err := db.Callback().Update().After("gorm:update").Register("test:auth-config-fail-second", func(tx *gorm.DB) {
		calls++
		tx.RowsAffected = 1
		if calls == 2 {
			tx.AddError(writeErr)
		}
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	repo := NewAuthConfigRepository(db)
	err := repo.UpdateValuesAtomically(context.Background(), map[string]string{
		"ZNS_TOKEN":   "new-access",
		"ZNS_REFRESH": "new-refresh",
	})
	if !errors.Is(err, writeErr) {
		t.Fatalf("UpdateValuesAtomically() error = %v, want wrapped %v", err, writeErr)
	}
	if calls != 2 {
		t.Fatalf("update calls = %d, want 2", calls)
	}
}

func TestUpdateValuesAtomicallyRejectsMissingConfigRow(t *testing.T) {
	db := newAuthConfigDryRunDB(t)
	if err := db.Callback().Update().After("gorm:update").Register("test:auth-config-no-row", func(tx *gorm.DB) {
		tx.RowsAffected = 0
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	repo := NewAuthConfigRepository(db)
	err := repo.UpdateValuesAtomically(context.Background(), map[string]string{
		"ZNS_TOKEN": "new-access",
	})
	if err == nil || !strings.Contains(err.Error(), "expected 1 row, got 0") {
		t.Fatalf("UpdateValuesAtomically() error = %v, want missing-row failure", err)
	}
}

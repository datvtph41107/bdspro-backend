package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"user/internal/subscriptioncutover"
)

func TestActivateRequiresExplicitGuard(t *testing.T) {
	t.Setenv("QHPRO_ALLOW_SUBSCRIPTION_READ_CUTOVER", "")
	err := run(context.Background(), options{action: "activate"})
	if !errors.Is(err, subscriptioncutover.ErrActivationDenied) {
		t.Fatalf("err = %v", err)
	}
}

func TestRollbackRequiresLegacyModeConfirmation(t *testing.T) {
	t.Setenv("QHPRO_ALLOW_SUBSCRIPTION_READ_ROLLBACK", "1")
	t.Setenv("QHPRO_CONFIRM_SUBSCRIPTION_READ_MODE_LEGACY", "")
	err := run(context.Background(), options{action: "rollback"})
	if err == nil {
		t.Fatal("rollback accepted without legacy mode confirmation")
	}
}

func TestParseTimeUsesUTC(t *testing.T) {
	value, err := parseTime("2026-08-04T13:00:00+07:00")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 8, 4, 6, 0, 0, 0, time.UTC)
	if !value.Equal(want) || value.Location() != time.UTC {
		t.Fatalf("value = %v", value)
	}
}

func TestWriteJSONIsAtomic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := writeJSON(path, map[string]any{"ready": true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temporary file remains: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "{\n  \"ready\": true\n}\n" {
		t.Fatalf("data=%q err=%v", data, err)
	}
}

package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"user/internal/modules/subscription/backfill"
)

func TestParseOptionalIDsDeduplicates(t *testing.T) {
	ids, err := parseOptionalIDs("9, 10,9")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 || ids[0] != 9 || ids[1] != 10 {
		t.Fatalf("ids = %v", ids)
	}
}

func TestParseOptionalIDsRejectsZero(t *testing.T) {
	if _, err := parseOptionalIDs("0"); err == nil {
		t.Fatal("zero membership id was accepted")
	}
}

func TestWriteJSONUsesRequestedPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.json")
	if err := writeJSON(path, map[string]int{"matched": 1}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "{\n  \"matched\": 1\n}\n" {
		t.Fatalf("content = %q", content)
	}
}

func TestRunRejectsExecuteWithoutExplicitGuard(t *testing.T) {
	t.Setenv("QHPRO_ALLOW_SUBSCRIPTION_BACKFILL_EXECUTE", "")
	err := run(context.Background(), options{action: "execute"})
	if !errors.Is(err, backfill.ErrExecutionForbidden) {
		t.Fatalf("err = %v", err)
	}
}

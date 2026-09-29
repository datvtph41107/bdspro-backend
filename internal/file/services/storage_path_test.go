package services

import (
	"path/filepath"
	"testing"
)

func TestValidatePathAllowsDescendant(t *testing.T) {
	root := filepath.Join("tmp", "files")
	child := filepath.Join(root, "public", "report.pdf")
	if err := validatePath(child, root); err != nil {
		t.Fatalf("expected descendant path to be valid: %v", err)
	}
}

func TestValidatePathRejectsEscape(t *testing.T) {
	root := filepath.Join("tmp", "files")
	escaped := filepath.Join(root, "..", "secret.txt")
	if err := validatePath(escaped, root); err == nil {
		t.Fatal("expected path escape to be rejected")
	}
}

func TestValidatePathAllowsDotsInsideFilename(t *testing.T) {
	root := filepath.Join("tmp", "files")
	child := filepath.Join(root, "public", "report..final.pdf")
	if err := validatePath(child, root); err != nil {
		t.Fatalf("dots inside a filename are not traversal: %v", err)
	}
}

func TestResolveStoredPathKeepsPathInsideStorageRoot(
	t *testing.T,
) {
	root := t.TempDir()

	storage, err := NewStorageService(root)
	if err != nil {
		t.Fatalf("NewStorageService() error = %v", err)
	}

	path, err := storage.ResolveStoredPath(
		filepath.Join(
			"public",
			"version",
			"42.apk",
		),
	)
	if err != nil {
		t.Fatalf(
			"ResolveStoredPath() error = %v",
			err,
		)
	}

	want := filepath.Join(
		root,
		"public",
		"version",
		"42.apk",
	)

	if path != want {
		t.Fatalf(
			"path = %q, want %q",
			path,
			want,
		)
	}
}

func TestResolveStoredPathRejectsStorageEscape(
	t *testing.T,
) {
	root := t.TempDir()

	storage, err := NewStorageService(root)
	if err != nil {
		t.Fatalf("NewStorageService() error = %v", err)
	}

	if _, err := storage.ResolveStoredPath(
		filepath.Join(
			"..",
			"outside.apk",
		),
	); err == nil {
		t.Fatal(
			"ResolveStoredPath() accepted path outside storage root",
		)
	}
}

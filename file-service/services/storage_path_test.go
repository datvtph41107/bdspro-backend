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

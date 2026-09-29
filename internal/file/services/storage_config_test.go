package services

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStorageRootIsPhysicalWhileStoredPathNamespaceIsStable(
	t *testing.T,
) {
	root := filepath.Join(t.TempDir(), "mounted-file-data")

	storage, err := NewStorageService(root)
	if err != nil {
		t.Fatalf("NewStorageService() error = %v", err)
	}

	if storage.PublicRoot() != filepath.Join(root, "public") {
		t.Fatalf("PublicRoot() = %q", storage.PublicRoot())
	}
	if storage.SecureRoot() != filepath.Join(root, "secure") {
		t.Fatalf("SecureRoot() = %q", storage.SecureRoot())
	}

	info, err := storage.StoreContent(
		strings.NewReader("payload"),
		"report.pdf",
		"application/pdf",
		int64(len("payload")),
		nil,
		42,
		"document",
	)
	if err != nil {
		t.Fatalf("StoreContent() error = %v", err)
	}

	if !strings.HasPrefix(
		filepath.ToSlash(info.RelativePath),
		"files/public/document/",
	) {
		t.Fatalf(
			"RelativePath = %q, want stable files/public namespace",
			info.RelativePath,
		)
	}
	if !strings.HasPrefix(
		filepath.Clean(info.AbsolutePath),
		filepath.Join(filepath.Clean(root), "public", "document")+
			string(filepath.Separator),
	) {
		t.Fatalf(
			"AbsolutePath = %q, want configured physical root",
			info.AbsolutePath,
		)
	}

	resolved, err := storage.ResolveStoredPath(info.RelativePath)
	if err != nil {
		t.Fatalf("ResolveStoredPath() error = %v", err)
	}
	if resolved != info.AbsolutePath {
		t.Fatalf(
			"ResolveStoredPath() = %q, want %q",
			resolved,
			info.AbsolutePath,
		)
	}
}

func TestNewStorageServiceFailsClosedForMissingRoot(t *testing.T) {
	if _, err := NewStorageService("   "); err == nil {
		t.Fatal("NewStorageService(empty) error = nil")
	}
}

func TestResolveStoredPathAcceptsCanonicalAndLegacyStorageRelativeInput(
	t *testing.T,
) {
	root := t.TempDir()
	storage, err := NewStorageService(root)
	if err != nil {
		t.Fatal(err)
	}

	canonical := filepath.ToSlash(
		filepath.Join("files", "public", "version", "42.apk"),
	)
	legacy := filepath.ToSlash(
		filepath.Join("public", "version", "42.apk"),
	)
	want := filepath.Join(root, "public", "version", "42.apk")

	for _, storedPath := range []string{canonical, legacy} {
		got, err := storage.ResolveStoredPath(storedPath)
		if err != nil {
			t.Fatalf(
				"ResolveStoredPath(%q) error = %v",
				storedPath,
				err,
			)
		}
		if got != want {
			t.Fatalf(
				"ResolveStoredPath(%q) = %q, want %q",
				storedPath,
				got,
				want,
			)
		}
	}
}

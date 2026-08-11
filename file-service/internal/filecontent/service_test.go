package filecontent

import (
	_utils "common/utils"
	"context"
	"errors"
	"file/services"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileRejectsPrefixThatDoesNotMatchStorageArea(t *testing.T) {
	root := t.TempDir()
	publicRoot := filepath.Join(root, "public")
	secureRoot := filepath.Join(root, "secure")
	if err := os.MkdirAll(secureRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	securePath := filepath.Join(secureRoot, "1.pdf")
	if err := os.WriteFile(securePath, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}

	const key = "test-key"
	service := NewService(nil, nil, &services.StorageService{
		RootLocation:   root,
		RootPublicPath: publicRoot,
		RootSecurePath: secureRoot,
	}, key)
	encodedAsPublic := "p" + _utils.XorEncode(filepath.Join("secure", "1.pdf"), key)

	_, err := service.ReadFile(context.Background(), encodedAsPublic, 0)
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("ReadFile() error = %v, want ErrInvalidPath", err)
	}
}

func TestReadSecureFileRequiresActor(t *testing.T) {
	root := t.TempDir()
	secureRoot := filepath.Join(root, "secure")
	if err := os.MkdirAll(secureRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(secureRoot, "1.pdf"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}

	const key = "test-key"
	service := NewService(nil, nil, &services.StorageService{
		RootLocation:   root,
		RootPublicPath: filepath.Join(root, "public"),
		RootSecurePath: secureRoot,
	}, key)
	encoded := "s" + _utils.XorEncode(filepath.Join("secure", "1.pdf"), key)

	_, err := service.ReadFile(context.Background(), encoded, 0)
	if !errors.Is(err, ErrNoActor) {
		t.Fatalf("ReadFile() error = %v, want ErrNoActor", err)
	}
}

func TestFileIDFromPath(t *testing.T) {
	fileID, err := fileIDFromPath("secure/document/2026-08-11/uuid/42.pdf")
	if err != nil || fileID != 42 {
		t.Fatalf("fileIDFromPath() = %d, %v", fileID, err)
	}
	if _, err := fileIDFromPath("secure/document/report.pdf"); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("invalid file ID error = %v", err)
	}
}

package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const storedPathNamespace = "files"

// validatePath verifies that a generated storage path stays under its canonical root.
// Invalid input is a request/data error, not a process-fatal programmer invariant.
func validatePath(path, parent string) error {
	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return fmt.Errorf("resolve storage path relative to root: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return fmt.Errorf("storage path escapes root")
	}
	return nil
}

// StoredPath projects a physical filesystem path into the stable File wire /
// durable namespace. The namespace is intentionally independent of the
// deployment-specific filesystem root.
func (s *StorageService) StoredPath(
	physicalPath string,
) (string, error) {
	if s == nil {
		return "", fmt.Errorf("storage is unavailable")
	}

	root := filepath.Clean(strings.TrimSpace(s.rootLocation))
	physicalPath = filepath.Clean(strings.TrimSpace(physicalPath))
	if root == "." || physicalPath == "." {
		return "", fmt.Errorf("storage path is empty")
	}
	if err := validatePath(physicalPath, root); err != nil {
		return "", err
	}

	rel, err := filepath.Rel(root, physicalPath)
	if err != nil {
		return "", fmt.Errorf("project storage path relative to root: %w", err)
	}
	if rel == "." || rel == "" {
		return "", fmt.Errorf("stored path must identify a file")
	}

	return filepath.ToSlash(
		filepath.Join(storedPathNamespace, rel),
	), nil
}

// ResolveStoredPath resolves a durable File wire path to its physical
// filesystem location and rejects paths outside the configured storage root.
//
// Canonical published paths start with "files/". Storage-relative paths
// without that prefix remain readable as a compatibility input for historical
// tests/data, but StoredPath always emits the canonical namespace.
func (s *StorageService) ResolveStoredPath(
	storedPath string,
) (string, error) {
	if s == nil {
		return "", fmt.Errorf("storage is unavailable")
	}

	storedPath = strings.TrimSpace(storedPath)
	if storedPath == "" {
		return "", fmt.Errorf("stored path is empty")
	}
	if filepath.IsAbs(storedPath) {
		return "", fmt.Errorf("stored path must be relative")
	}

	cleanSlash := filepath.ToSlash(filepath.Clean(storedPath))
	if cleanSlash == storedPathNamespace {
		return "", fmt.Errorf("stored path must identify a file")
	}
	if strings.HasPrefix(cleanSlash, storedPathNamespace+"/") {
		cleanSlash = strings.TrimPrefix(
			cleanSlash,
			storedPathNamespace+"/",
		)
	}

	root := filepath.Clean(s.rootLocation)
	physicalPath := filepath.Join(
		root,
		filepath.FromSlash(cleanSlash),
	)

	if err := validatePath(physicalPath, root); err != nil {
		return "", err
	}

	return physicalPath, nil
}

// RemovePhysicalPath removes one File-owned physical artifact while enforcing
// the configured storage root boundary. It is used only for compensation of a
// failed durable effect; durable/public paths are never accepted here.
func (s *StorageService) RemovePhysicalPath(
	physicalPath string,
) error {
	if s == nil {
		return fmt.Errorf("storage is unavailable")
	}

	root := filepath.Clean(strings.TrimSpace(s.rootLocation))
	physicalPath = filepath.Clean(strings.TrimSpace(physicalPath))
	if root == "." || physicalPath == "." {
		return fmt.Errorf("storage path is empty")
	}
	if physicalPath == root {
		return fmt.Errorf("refusing to remove storage root")
	}
	if err := validatePath(physicalPath, root); err != nil {
		return err
	}

	if err := os.RemoveAll(physicalPath); err != nil {
		return fmt.Errorf("remove physical storage path: %w", err)
	}
	return nil
}

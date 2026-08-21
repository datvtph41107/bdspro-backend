package services

import (
	"fmt"
	"path/filepath"
	"strings"
)

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

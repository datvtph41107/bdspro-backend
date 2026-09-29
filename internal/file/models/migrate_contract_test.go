package models

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonicalSQLCoversFileAutoMigrateRegistry(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob(filepath.Join("..", "database", "migrations", "*.up.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, files)
	created := map[string]struct{}{}
	pattern := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-z0-9_]+)`)
	for _, file := range files {
		contents, readErr := os.ReadFile(file)
		require.NoError(t, readErr)
		for _, match := range pattern.FindAllStringSubmatch(string(contents), -1) {
			created[strings.ToLower(match[1])] = struct{}{}
		}
	}

	require.Contains(t, created, "file")
	require.Contains(t, created, "file_access")
}

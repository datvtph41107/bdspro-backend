package postgres

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm/schema"
)

var (
	paymentCreateTablePattern = regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?"?([a-z0-9_]+)"?\s*\((.*?)\);`)
	paymentColumnPattern      = regexp.MustCompile(`(?m)^\s*"?([a-z_][a-z0-9_]*)"?\s+`)
)

func TestCanonicalSQLCoversPaymentAutoMigrateRegistry(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob(filepath.Join("..", "..", "database", "migrations", "*.up.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	tables := make(map[string]map[string]struct{})
	for _, filename := range files {
		contents, readErr := os.ReadFile(filename)
		require.NoError(t, readErr)
		for _, tableMatch := range paymentCreateTablePattern.FindAllStringSubmatch(string(contents), -1) {
			columns := make(map[string]struct{})
			for _, columnMatch := range paymentColumnPattern.FindAllStringSubmatch(tableMatch[2], -1) {
				columns[strings.ToLower(columnMatch[1])] = struct{}{}
			}
			tables[strings.ToLower(tableMatch[1])] = columns
		}
	}

	cache := &sync.Map{}
	for _, model := range AutoMigrateModels() {
		modelSchema, parseErr := schema.Parse(model, cache, schema.NamingStrategy{})
		require.NoError(t, parseErr)
		columns, found := tables[strings.ToLower(modelSchema.Table)]
		require.True(t, found, "canonical Payment SQL is missing GORM table %s", modelSchema.Table)
		for _, field := range modelSchema.Fields {
			if field.DBName != "" {
				require.Contains(t, columns, strings.ToLower(field.DBName),
					"canonical Payment SQL table %s is missing GORM column %s", modelSchema.Table, field.DBName)
			}
		}
	}
}

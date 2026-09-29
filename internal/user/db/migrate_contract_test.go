package db

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
	createTablePattern = regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?"?([a-z0-9_]+)"?\s*\((.*?)\);`)
	columnPattern      = regexp.MustCompile(`(?m)^\s*"?([a-z_][a-z0-9_]*)"?\s+`)
	dropTablePattern   = regexp.MustCompile(`(?i)DROP\s+TABLE\s+(?:IF\s+EXISTS\s+)?"?([a-z0-9_]+)"?`)
)

func TestCanonicalSQLCoversUserAutoMigrateRegistry(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile(filepath.Join("..", "database", "migrations", "000001_user_schema.up.sql"))
	require.NoError(t, err)
	tables := sqlTables(string(contents))

	// Parse đúng metadata mà GORM dùng thay vì duy trì một danh sách table thủ
	// công có thể lệch khỏi registry AutoMigrate.
	cache := &sync.Map{}
	for _, model := range AutoMigrateModels() {
		modelSchema, parseErr := schema.Parse(model, cache, schema.NamingStrategy{})
		require.NoError(t, parseErr)
		requireSchemaCovered(t, tables, modelSchema.Table, modelSchema.Fields)

		for _, relationship := range modelSchema.Relationships.Relations {
			if relationship.JoinTable == nil {
				continue
			}
			requireSchemaCovered(t, tables, relationship.JoinTable.Table, relationship.JoinTable.Fields)
		}
	}
}

func TestCanonicalUserDownCoversEveryCreatedTable(t *testing.T) {
	t.Parallel()

	upContents, err := os.ReadFile(filepath.Join("..", "database", "migrations", "000001_user_schema.up.sql"))
	require.NoError(t, err)
	downContents, err := os.ReadFile(filepath.Join("..", "database", "migrations", "000001_user_schema.down.sql"))
	require.NoError(t, err)

	dropped := make(map[string]struct{})
	for _, match := range dropTablePattern.FindAllStringSubmatch(string(downContents), -1) {
		dropped[strings.ToLower(match[1])] = struct{}{}
	}
	for table := range sqlTables(string(upContents)) {
		require.Contains(t, dropped, table, "User down migration is missing table %s", table)
	}
}

func sqlTables(contents string) map[string]map[string]struct{} {
	tables := make(map[string]map[string]struct{})
	for _, tableMatch := range createTablePattern.FindAllStringSubmatch(contents, -1) {
		table := strings.ToLower(tableMatch[1])
		columns := make(map[string]struct{})
		for _, columnMatch := range columnPattern.FindAllStringSubmatch(tableMatch[2], -1) {
			columns[strings.ToLower(columnMatch[1])] = struct{}{}
		}
		tables[table] = columns
	}
	return tables
}

func requireSchemaCovered(t *testing.T, tables map[string]map[string]struct{}, table string, fields []*schema.Field) {
	t.Helper()
	table = strings.ToLower(table)
	columns, found := tables[table]
	require.True(t, found, "canonical User SQL is missing GORM table %s", table)
	for _, field := range fields {
		if field.DBName == "" {
			continue
		}
		require.Contains(t, columns, strings.ToLower(field.DBName),
			"canonical User SQL table %s is missing GORM column %s", table, field.DBName)
	}
}

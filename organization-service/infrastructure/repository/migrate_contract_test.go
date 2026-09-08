package repository

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
	organizationCreateTablePattern = regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?"?([a-z0-9_]+)"?\s*\((.*?)\);`)
	organizationColumnPattern      = regexp.MustCompile(`(?im)(?:^|,)\s*"?([a-z_][a-z0-9_]*)"?\s+`)
	organizationDropTablePattern   = regexp.MustCompile(`(?i)DROP\s+TABLE\s+(?:IF\s+EXISTS\s+)?"?([a-z0-9_]+)"?`)
)

func TestCanonicalSQLCoversOrganizationAutoMigrateRegistry(t *testing.T) {
	t.Parallel()

	up, err := os.ReadFile(filepath.Join("..", "..", "migrate", "000001_organization_schema.up.sql"))
	require.NoError(t, err)
	tables := organizationSQLTables(string(up))

	cache := &sync.Map{}
	for _, model := range AutoMigrateModels() {
		modelSchema, parseErr := schema.Parse(model, cache, schema.NamingStrategy{})
		require.NoError(t, parseErr)
		requireOrganizationSchemaCovered(t, tables, modelSchema.Table, modelSchema.Fields)
		for _, relationship := range modelSchema.Relationships.Relations {
			if relationship.JoinTable != nil {
				requireOrganizationSchemaCovered(t, tables, relationship.JoinTable.Table, relationship.JoinTable.Fields)
			}
		}
	}
}

func TestCanonicalOrganizationDownCoversEveryCreatedTable(t *testing.T) {
	t.Parallel()

	up, err := os.ReadFile(filepath.Join("..", "..", "migrate", "000001_organization_schema.up.sql"))
	require.NoError(t, err)
	down, err := os.ReadFile(filepath.Join("..", "..", "migrate", "000001_organization_schema.down.sql"))
	require.NoError(t, err)

	dropped := make(map[string]struct{})
	for _, match := range organizationDropTablePattern.FindAllStringSubmatch(string(down), -1) {
		dropped[strings.ToLower(match[1])] = struct{}{}
	}
	for table := range organizationSQLTables(string(up)) {
		require.Contains(t, dropped, table, "Organization down migration is missing table %s", table)
	}
}

func organizationSQLTables(contents string) map[string]map[string]struct{} {
	tables := make(map[string]map[string]struct{})
	for _, tableMatch := range organizationCreateTablePattern.FindAllStringSubmatch(contents, -1) {
		columns := make(map[string]struct{})
		for _, columnMatch := range organizationColumnPattern.FindAllStringSubmatch(tableMatch[2], -1) {
			columns[strings.ToLower(columnMatch[1])] = struct{}{}
		}
		tables[strings.ToLower(tableMatch[1])] = columns
	}
	return tables
}

func requireOrganizationSchemaCovered(
	t *testing.T,
	tables map[string]map[string]struct{},
	table string,
	fields []*schema.Field,
) {
	t.Helper()
	columns, found := tables[strings.ToLower(table)]
	require.True(t, found, "canonical Organization SQL is missing GORM table %s", table)
	for _, field := range fields {
		if field.DBName != "" {
			require.Contains(t, columns, strings.ToLower(field.DBName),
				"canonical Organization SQL table %s is missing GORM column %s", table, field.DBName)
		}
	}
}

package db

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var createTablePattern = regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-z0-9_]+)`)
var migrationFilePattern = regexp.MustCompile(`^([0-9]{6})_.+\.(up|down)\.sql$`)
var referencePattern = regexp.MustCompile(`(?i)REFERENCES\s+(?:[a-z0-9_]+\.)?"?([a-z0-9_]+)"?`)
var lineCommentPattern = regexp.MustCompile(`(?m)--[^\n]*$`)

func TestCanonicalMigrationSequence(t *testing.T) {
	t.Parallel()
	const latestMigration = 44

	entries, err := os.ReadDir(filepath.Join("..", "..", "migrate"))
	require.NoError(t, err)

	pairs := make(map[string]map[string]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		match := migrationFilePattern.FindStringSubmatch(entry.Name())
		require.Len(t, match, 3, "invalid migration filename %s", entry.Name())
		version, direction := match[1], match[2]
		if pairs[version] == nil {
			pairs[version] = make(map[string]string)
		}
		require.Empty(t, pairs[version][direction], "duplicate %s migration for version %s", direction, version)
		pairs[version][direction] = entry.Name()
	}

	require.Len(t, pairs, latestMigration)
	for version := 1; version <= latestMigration; version++ {
		key := fmt.Sprintf("%06d", version)
		require.Contains(t, pairs, key, "migration history has a version gap")
		require.NotEmpty(t, pairs[key]["up"], "version %s has no up migration", key)
		require.NotEmpty(t, pairs[key]["down"], "version %s has no down migration", key)
	}
}

func TestCanonicalSQLCoversRuntimeSchemaRegistry(t *testing.T) {
	t.Parallel()

	migrationDir := filepath.Join("..", "..", "migrate")
	entries, err := os.ReadDir(migrationDir)
	require.NoError(t, err)

	created := make(map[string]struct{})
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		contents, readErr := os.ReadFile(filepath.Join(migrationDir, entry.Name()))
		require.NoError(t, readErr)
		for _, match := range createTablePattern.FindAllStringSubmatch(string(contents), -1) {
			created[strings.ToLower(match[1])] = struct{}{}
		}
	}

	requiredTables := []string{
		"amenities",
		"contact_labels",
		"directory_categories",
		"directory_sources",
		"directory_suppliers",
		"features",
		"map_points",
		"open_hours",
		"parcels",
		"poi_amenities",
		"poi_categories",
		"poi_media",
		"pois",
		"province_v2",
		"pro_ai_jobs",
		"qh_audit_entries",
		"qh_authority_issuring",
		"qh_direction",
		"qh_jurisdictions",
		"qh_label_layers",
		"qh_labels",
		"qh_land_use",
		"qh_layer_families",
		"qh_layer_land_use",
		"qh_layer_legals",
		"qh_layer_resolver_configs",
		"qh_layers",
		"qh_legal_documents",
		"qh_legends",
		"qh_parcel_direction",
		"qh_parcel_info",
		"qh_planning_documents",
		"qh_planning_events",
		"qh_planning_land_use",
		"qh_planning_projects",
		"qh_planning_relations",
		"qh_region_extends",
		"qh_region_import_error_logs",
		"qh_regions",
		"qh_shape",
		"quota_usage_events",
		"report_events",
		"report_jobs",
		"reports",
		"search_index",
		"user_followed_parcels",
		"user_followed_planning_projects",
		"user_notifications",
		"user_reported",
		"user_subscriptions",
		"user_view_events",
		"user_view_history",
		"ward_v2",
	}

	missing := make([]string, 0)
	for _, table := range requiredTables {
		if _, ok := created[table]; !ok {
			missing = append(missing, table)
		}
	}
	sort.Strings(missing)
	require.Empty(t, missing, "canonical SQL does not cover the runtime schema registry")
}

func TestEveryForeignKeyTargetExistsBeforeFirstUse(t *testing.T) {
	t.Parallel()

	migrationDir := filepath.Join("..", "..", "migrate")
	entries, err := os.ReadDir(migrationDir)
	require.NoError(t, err)
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// Duyệt theo đúng thứ tự version và vị trí câu lệnh trong file. Cách này bắt
	// được cả dependency khác version lẫn FOREIGN KEY nằm trước CREATE TABLE
	// prerequisite trong cùng một baseline.
	created := make(map[string]struct{})
	for _, filename := range files {
		contents, readErr := os.ReadFile(filepath.Join(migrationDir, filename))
		require.NoError(t, readErr)
		contents = lineCommentPattern.ReplaceAll(contents, nil)
		type schemaEvent struct {
			offset    int
			table     string
			reference bool
		}
		events := make([]schemaEvent, 0)
		for _, location := range createTablePattern.FindAllSubmatchIndex(contents, -1) {
			events = append(events, schemaEvent{offset: location[0], table: strings.ToLower(string(contents[location[2]:location[3]]))})
		}
		for _, location := range referencePattern.FindAllSubmatchIndex(contents, -1) {
			events = append(events, schemaEvent{offset: location[0], table: strings.ToLower(string(contents[location[2]:location[3]])), reference: true})
		}
		sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
		for _, event := range events {
			if !event.reference {
				created[event.table] = struct{}{}
				continue
			}
			_, found := created[event.table]
			require.True(t, found, "%s references table %s before it is created", filename, event.table)
		}
	}
}

func TestSchemaCompletionRollbackIsNonDestructive(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile(filepath.Join("..", "..", "migrate", "000041_complete_gorm_owned_schema.down.sql"))
	require.NoError(t, err)
	sql := strings.ToLower(string(contents))
	require.NotContains(t, sql, "drop table")
	require.NotContains(t, sql, "cascade")
}

func TestLayerFamilySortOrderIsVersioned(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile(filepath.Join("..", "..", "migrate", "000043_add_qh_layer_family_sort_number.up.sql"))
	require.NoError(t, err)
	sql := strings.ToLower(string(contents))
	require.Contains(t, sql, "add column if not exists sort_number integer not null default 0")
	require.Contains(t, sql, "idx_qh_layer_families_active_sort")
}

func TestAdministrativeCatalogCutoverIsVersioned(t *testing.T) {
	t.Parallel()

	contents, err := os.ReadFile(filepath.Join("..", "..", "migrate", "000044_canonicalize_vn_administrative_catalog.up.sql"))
	if err != nil {
		t.Fatalf("read administrative catalog migration: %v", err)
	}
	sql := string(contents)
	for _, fragment := range []string{
		"INSERT INTO province_v2",
		"FROM provinces",
		"INSERT INTO ward_v2",
		"FROM wards AS legacy_ward",
		"ward_count = counts.ward_count",
		"ON CONFLICT DO NOTHING",
	} {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("administrative catalog migration missing %q", fragment)
		}
	}
}

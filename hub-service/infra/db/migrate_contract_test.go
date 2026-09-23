package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCanonicalSQLCoversHubAutoMigrateRegistry(t *testing.T) {
	entries, err := os.ReadDir("../../migrate")
	if err != nil {
		t.Fatal(err)
	}

	var sqlBuilder strings.Builder
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		body, err := os.ReadFile(filepath.Join("../../migrate", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		sqlBuilder.Write(body)
		sqlBuilder.WriteByte('\n')
	}
	source := strings.ToLower(sqlBuilder.String())

	for _, table := range []string{
		"applink", "tb_event_queue", "provinces", "districts", "wards",
		"province_v2", "ward_v2", "tb_user_guide", "system_config", "faqs",
		"versions", "api_keys", "user_guide_steps", "interactive_events",
		"error_logs", "update_data",
		"click_urls", "request_stats_log", "user_interaction_log",
		"request_audit_log", "keyword_search_log", "search_latency_log",
	} {
		if !strings.Contains(source, "create table "+table) &&
			!strings.Contains(source, "create table if not exists "+table) {
			t.Errorf("canonical Hub migration thiếu table %s", table)
		}
	}
	if got := len(AutoMigrateModels()); got != 22 {
		t.Fatalf("AutoMigrateModels() có %d model, muốn 22", got)
	}
}

func TestListPendingMigrationsSkipsApplied(t *testing.T) {
	dir := t.TempDir()
	mustWrite := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("000003_hub_fix_planning_faq_link.up.sql", "SELECT 1;")
	mustWrite("000004_hub_request_stats_audit.up.sql", "SELECT 1;")
	mustWrite("readme.txt", "ignore")

	pending, err := listPendingMigrations(dir, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 1 || pending[0].version != 4 {
		t.Fatalf("pending = %#v, want only version 4", pending)
	}
}

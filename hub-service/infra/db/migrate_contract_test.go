package db

import (
	"os"
	"strings"
	"testing"
)

func TestCanonicalSQLCoversHubAutoMigrateRegistry(t *testing.T) {
	sql, err := os.ReadFile("../../migrate/000001_hub_schema.up.sql")
	if err != nil {
		t.Fatal(err)
	}

	source := strings.ToLower(string(sql))
	for _, table := range []string{
		"applink", "tb_event_queue", "provinces", "districts", "wards",
		"province_v2", "ward_v2", "tb_user_guide", "system_config", "faqs",
		"versions", "api_keys", "user_guide_steps", "interactive_events",
		"error_logs", "update_data",
	} {
		if !strings.Contains(source, "create table "+table) {
			t.Errorf("canonical Hub migration thiếu table %s", table)
		}
	}
	if got := len(AutoMigrateModels()); got != 16 {
		t.Fatalf("AutoMigrateModels() có %d model, muốn 16", got)
	}
}

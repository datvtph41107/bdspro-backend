package db

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonicalSQLCoversNotificationAutoMigrateRegistry(t *testing.T) {
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

	required := []string{
		"account_warning_logs", "account_warning_templates", "account_warnings", "admin_histories", "deal_history",
		"event_notifications", "history_auth", "noti_histories", "notification", "notification_delivery_intents",
		"notification_history", "payment_event_inbox", "person_configs", "property_histories", "tb_history",
	}
	for _, table := range required {
		require.Contains(t, created, table, "canonical Notification SQL is missing runtime table %s", table)
	}
}

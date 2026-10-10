package postgres

import "strings"

func firstNonEmptyDSN(values ...string) string {
	for _, value := range values {
		if dsn := strings.TrimSpace(value); dsn != "" {
			return dsn
		}
	}
	return ""
}

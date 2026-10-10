package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestFormatDevSQLLine(t *testing.T) {
	line := formatDevSQLLine(
		time.Date(2026, 10, 9, 8, 13, 58, 0, time.FixedZone("+7", 7*3600)),
		slog.LevelWarn,
		map[string]string{
			"duration_ms":  "1334.082",
			"rows":         "0",
			"sql":          "SELECT DISTINCT pp.id\n\t\tFROM parcels p",
			"query_source": "/Users/baphieu/Workspace/ME/go-microservices/v4/modules/tqd/postgres/publiccontent/repository.go:393",
		},
	)
	if !strings.Contains(line, "1334.082ms") {
		t.Fatalf("unexpected sql line: %s", line)
	}
	if !strings.Contains(line, "repository.go:393") {
		t.Fatalf("expected short source, got: %s", line)
	}
	if !strings.Contains(line, "SELECT DISTINCT") || !strings.Contains(line, "FROM parcels p") {
		t.Fatalf("expected full sql on second line: %s", line)
	}
}

func TestDevTailAttrsNoDuplicatesKeepsHTTPPath(t *testing.T) {
	got := devTailAttrs(map[string]string{
		"category":         "application",
		"error":            "record not found",
		"http_status_code": "500",
		"http_method":      "POST",
		"http_path":        "/v2/auth/token/refresh",
		"profile_id":       "267",
	})
	want := "error=record not found http_status_code=500 http_method=POST http_path=/v2/auth/token/refresh profile_id=267"
	if got != want {
		t.Fatalf("got  %q\nwant %q", got, want)
	}
}

func TestFormatDevSQLLineIncludesError(t *testing.T) {
	line := formatDevSQLLine(
		time.Date(2026, 10, 9, 9, 45, 30, 0, time.UTC),
		slog.LevelError,
		map[string]string{
			"duration_ms": "52.559",
			"rows":        "0",
			"sql":         "INSERT INTO seo_count_log (seo_id, clicked, \"timestamp\") VALUES ($1, 1, $2)",
			"error":       "ERROR: relation \"seo_count_log\" does not exist (SQLSTATE 42P01)",
			"query_source": "news.go:467",
		},
	)
	if !strings.Contains(line, "INSERT INTO seo_count_log") {
		t.Fatalf("expected sql in line: %s", line)
	}
	if !strings.Contains(line, "relation \"seo_count_log\" does not exist") {
		t.Fatalf("expected postgres error in line: %s", line)
	}
}

func TestFormatLatencyDisplay(t *testing.T) {
	if got := formatLatencyDisplay(150 * time.Millisecond); got != "  150.00ms" {
		t.Fatalf("sub-second latency: %q", got)
	}
	if got := formatLatencyDisplay(2500 * time.Millisecond); got != "     2500ms" {
		t.Fatalf("multi-second latency: %q", got)
	}
}

func TestDevConsoleHandlerHTTPAccessMessage(t *testing.T) {
	var output bytes.Buffer
	c := defaults()
	c.ConsoleFormat = "dev"
	h := newHandler(c, &output, "")
	log := slog.New(h)
	log.Info("GET /v2/tqd/public/parcels/536061/quick-view 3.211514208s", "category", "application")
	got := output.String()
	if !strings.Contains(got, "[GIN]") || !strings.Contains(got, "GET") || !strings.Contains(got, "/v2/tqd/public/parcels/536061/quick-view") {
		t.Fatalf("unexpected access line: %s", got)
	}
}

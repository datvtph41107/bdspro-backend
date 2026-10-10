package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"
)

type devConsoleHandler struct {
	out  io.Writer
	opts slog.HandlerOptions
	attrs []slog.Attr
	groups []string
}

func newDevConsoleHandler(out io.Writer, opts *slog.HandlerOptions) *devConsoleHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	return &devConsoleHandler{out: out, opts: *opts}
}

func (h *devConsoleHandler) Enabled(_ context.Context, level slog.Level) bool {
	min := slog.LevelInfo
	if h.opts.Level != nil {
		min = h.opts.Level.Level()
	}
	return level >= min
}

func (h *devConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := map[string]string{}
	for _, a := range h.attrs {
		attrs[a.Key] = a.Value.String()
	}
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.String()
		return true
	})
	line := formatDevConsoleLine(r.Time, r.Level, r.Message, attrs)
	if line == "" {
		return nil
	}
	_, err := fmt.Fprintln(h.out, line)
	return err
}

func (h *devConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = append(append([]slog.Attr(nil), h.attrs...), attrs...)
	return &next
}

func (h *devConsoleHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	next := *h
	next.groups = append(append([]string(nil), h.groups...), name)
	return &next
}

func DevConsoleEnabled() bool {
	c := currentConfig.Load()
	return c != nil && c.ConsoleFormat == "dev"
}

func formatDevConsoleLine(ts time.Time, level slog.Level, msg string, attrs map[string]string) string {
	switch attrs["event_name"] {
	case "database.query.completed":
		return formatDevSQLLine(ts, level, attrs)
	case "process.logging.ready", "process.logging.stopped":
		return fmt.Sprintf("%s %-5s %s", ts.Format("15:04:05"), colorLevel(levelString(level)), msg)
	}

	if category := attrs["category"]; category == "http" || looksLikeHTTPAccess(msg) {
		return formatDevHTTPFromMessage(ts, msg, attrs)
	}

	lvl := colorLevel(levelString(level))
	if level >= slog.LevelError {
		return fmt.Sprintf("%s %-5s %s %s", ts.Format("15:04:05"), lvl, msg, devTailAttrs(attrs))
	}
	if level >= slog.LevelWarn && attrs["category"] == "database" {
		return formatDevSQLLine(ts, level, attrs)
	}

	extra := devTailAttrs(attrs)
	if extra == "" {
		return fmt.Sprintf("%s %-5s %s", ts.Format("15:04:05"), lvl, msg)
	}
	return fmt.Sprintf("%s %-5s %s %s", ts.Format("15:04:05"), lvl, msg, extra)
}

func formatDevSQLLine(ts time.Time, level slog.Level, attrs map[string]string) string {
	source := shortSource(attrs["query_source"])
	header := fmt.Sprintf(
		"%s %-5s SQL %s rows=%s %s",
		ts.Format("15:04:05"),
		colorLevel(levelString(level)),
		formatSQLDuration(attrs["duration_ms"]),
		attrs["rows"],
		source,
	)
	var b strings.Builder
	b.WriteString(header)
	if sqlText := normalizeSQL(attrs["sql"]); sqlText != "" {
		b.WriteString("\n        ")
		b.WriteString(sqlText)
	}
	if errText := strings.TrimSpace(attrs["error"]); errText != "" {
		b.WriteString("\n        ")
		b.WriteString(errText)
	}
	if b.Len() == len(header) {
		return header
	}
	return b.String()
}

func formatSQLDuration(durationMS string) string {
	durationMS = strings.TrimSpace(durationMS)
	if durationMS == "" {
		return "-"
	}
	if strings.HasSuffix(durationMS, "ms") {
		return durationMS
	}
	return durationMS + "ms"
}

func formatDevHTTPFromMessage(ts time.Time, msg string, attrs map[string]string) string {
	method, path, latency := parseHTTPAccessMessage(msg)
	if method == "" {
		method = attrs["http_method"]
	}
	if path == "" {
		path = attrs["http_path"]
	}
	status := attrs["http_status_code"]
	if status == "" {
		status = "200"
	}
	clientIP := attrs["client_ip"]
	if clientIP == "" {
		clientIP = "-"
	}
	latencyDur := parseDurationAttr(latency, attrs["duration"])
	statusCode := 200
	if status != "" {
		if _, err := fmt.Sscanf(status, "%d", &statusCode); err != nil {
			statusCode = 200
		}
	}
	return FormatGinAccessLine(ts, statusCode, latencyDur, clientIP, method, path)
}

func parseDurationAttr(primary, fallback string) time.Duration {
	for _, raw := range []string{primary, fallback} {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if d, err := time.ParseDuration(raw); err == nil {
			return d
		}
	}
	return 0
}

func FormatGinAccessLine(ts time.Time, statusCode int, latency time.Duration, clientIP, method, path string) string {
	if method == "" {
		method = "-"
	}
	if path == "" {
		path = "/"
	}
	if clientIP == "" {
		clientIP = "-"
	}
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	latencyText := formatLatencyDisplay(latency)
	if !devConsoleColorEnabled() {
		return fmt.Sprintf(
			"[GIN] %s | %3d |%s| %-15s | %-7s %s",
			ts.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latencyText,
			clientIP,
			method,
			path,
		)
	}

	statusColor := statusCodeColor(statusCode)
	methodCol := methodColor(method)
	latencyCol := latencyColor(latency)
	return fmt.Sprintf(
		"[GIN] %s |%s %3d %s|%s%s%s| %-15s |%s %-7s %s %s",
		ts.Format("2006/01/02 - 15:04:05"),
		statusColor, statusCode, ansiReset,
		latencyCol, latencyText, ansiReset,
		clientIP,
		methodCol, method, ansiReset,
		path,
	)
}

func looksLikeHTTPAccess(msg string) bool {
	fields := strings.Fields(msg)
	if len(fields) < 3 {
		return false
	}
	switch fields[0] {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
		return strings.HasPrefix(fields[1], "/")
	default:
		return false
	}
}

func parseHTTPAccessMessage(msg string) (method, path, latency string) {
	fields := strings.Fields(msg)
	if len(fields) < 3 {
		return "", "", ""
	}
	method = fields[0]
	path = fields[1]
	latency = fields[len(fields)-1]
	return method, path, latency
}

func normalizeSQL(sql string) string {
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	return strings.Join(strings.Fields(sql), " ")
}

func shortSource(source string) string {
	if source == "" {
		return ""
	}
	if idx := strings.LastIndex(source, "/"); idx >= 0 {
		source = source[idx+1:]
	}
	return source
}

func devTailAttrs(attrs map[string]string) string {
	skip := map[string]struct{}{
		"category": {}, "event_name": {}, "service_name": {}, "environment": {}, "run_id": {},
		"sql": {}, "rows": {}, "duration_ms": {}, "query_source": {}, "db_system": {}, "sql_operation": {},
	}
	const maxExtra = 4
	parts := make([]string, 0, 8)
	for _, key := range []string{"error", "http_status_code", "http_method", "http_path", "profile_id", "auth_id", "request_id", "legacy_source"} {
		if v := strings.TrimSpace(attrs[key]); v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", key, v))
		}
		skip[key] = struct{}{}
	}
	rest := make([]string, 0, len(attrs))
	for key, value := range attrs {
		if _, ok := skip[key]; ok || strings.TrimSpace(value) == "" {
			continue
		}
		rest = append(rest, key)
	}
	sort.Strings(rest)
	if len(rest) > maxExtra {
		rest = rest[:maxExtra]
	}
	for _, key := range rest {
		parts = append(parts, fmt.Sprintf("%s=%s", key, attrs[key]))
	}
	return strings.Join(parts, " ")
}

func levelString(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARN"
	case level >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}

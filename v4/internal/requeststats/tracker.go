package requeststats

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	_jwt "common/jwt"
)

var (
	numericPathSegment = regexp.MustCompile(`^[0-9]+$`)
	uuidPathSegment    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// Tracker observes HTTP requests and forwards them to a local Collector.
type Tracker struct {
	collector      *Collector
	specURLs       map[string]struct{}
	ignorePrefixes []string
}

// New builds a tracker that writes into the process-owned database via collector.
func New(collector *Collector, specURLs, ignorePrefixes []string) *Tracker {
	set := make(map[string]struct{}, len(specURLs))
	for _, raw := range specURLs {
		path := strings.TrimSpace(raw)
		if path == "" {
			continue
		}
		set[path] = struct{}{}
	}
	return &Tracker{
		collector:      collector,
		specURLs:       set,
		ignorePrefixes: mergeIgnorePrefixes(ignorePrefixes),
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

// Middleware records each HTTP request into local request_stats_log /
// user_interaction_log / request_audit_log (buffered flush).
// Mount after JWT auth so profile/organization ids are available when present.
func Middleware(tracker *Tracker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: 0}
			next.ServeHTTP(rec, r)

			if tracker == nil || tracker.collector == nil {
				return
			}

			path := r.URL.Path
			if tracker.shouldSkip(path) {
				return
			}

			normalizedPath := normalizePath(path)
			isSpecAudit := tracker.isSpecURL(path) || tracker.isSpecURL(normalizedPath)

			var profileID uint64
			var organizationID uint64
			if principal, ok := _jwt.PrincipalFromRequestContext(r.Context()); ok && principal.ProfileId > 0 {
				profileID = principal.ProfileId
				if principal.OrganizationId != nil {
					organizationID = *principal.OrganizationId
				}
			}

			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}

			tracker.collector.Track(TrackEvent{
				Method:         r.Method,
				Path:           normalizedPath,
				Service:        extractService(normalizedPath),
				StatusCode:     int32(status),
				DurationMs:     time.Since(start).Milliseconds(),
				ProfileID:      profileID,
				OrganizationID: organizationID,
				IsSpecAudit:    isSpecAudit,
			})
		})
	}
}

var defaultIgnorePrefixes = []string{
	"/swagger",
	"/swagger-doc",
	"/health",
	"/healthz",
	"/readyz",
	"/v2/hub/request-stats",
	"/v2/hub/keyword-search-stats",
	"/v2/hub/search-latency",
	"/v2/crm/public/seo-domains/top-clicks",
	"/v2/crm/public/seo-domains/click-daily",
	"/v2/user/admin/users/dashboard",
	"/v2/user/stats",
	"/v2/payment/wallets/dashboard",
	"/v2/payment/dashboard-metrics",
	"/v2/payment/stats",
	"/v2/bdspro/dashboard",
	"/v2/bdspro/v2/dashboard",
	"/v2/org/organization/current/dashboard",
	"/v2/admin-history/search",
}

func mergeIgnorePrefixes(extra []string) []string {
	seen := make(map[string]struct{}, len(defaultIgnorePrefixes)+len(extra))
	out := make([]string, 0, len(defaultIgnorePrefixes)+len(extra))
	add := func(raw string) {
		prefix := strings.TrimSpace(raw)
		if prefix == "" {
			return
		}
		if _, ok := seen[prefix]; ok {
			return
		}
		seen[prefix] = struct{}{}
		out = append(out, prefix)
	}
	for _, prefix := range defaultIgnorePrefixes {
		add(prefix)
	}
	for _, prefix := range extra {
		add(prefix)
	}
	return out
}

func (t *Tracker) shouldSkip(path string) bool {
	if path == "" {
		return true
	}
	prefixes := defaultIgnorePrefixes
	if t != nil && len(t.ignorePrefixes) > 0 {
		prefixes = t.ignorePrefixes
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func (t *Tracker) isSpecURL(path string) bool {
	if t == nil || len(t.specURLs) == 0 {
		return false
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	_, ok := t.specURLs[path]
	return ok
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/"
	}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "" {
			continue
		}
		if isPathIDSegment(part) {
			parts[i] = ":id"
		}
	}
	normalized := strings.Join(parts, "/")
	if normalized == "" {
		return "/"
	}
	return normalized
}

func isPathIDSegment(part string) bool {
	switch strings.ToLower(strings.TrimSpace(part)) {
	case "undefined", "null", "nan", "nil":
		return true
	}
	return numericPathSegment.MatchString(part) || uuidPathSegment.MatchString(part)
}

func extractService(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return "unknown"
	}
	version := parts[0]
	if version != "v1" && version != "v2" && version != "v2.1" && version != "v3" {
		return "unknown"
	}
	service := strings.TrimSpace(parts[1])
	if service == "" || service == ":id" {
		return "unknown"
	}
	return service
}

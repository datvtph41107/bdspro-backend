package requeststats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	_httpauth "gateway/internal/httpauth"

	hubpb "pb/types/hub"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	numericPathSegment = regexp.MustCompile(`^[0-9]+$`)
	uuidPathSegment    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// Tracker gửi thống kê request (fire-and-forget).
// Ưu tiên TrackURL → v4 HTTP; fallback hub gRPC khi TrackURL trống.
type Tracker struct {
	client         hubpb.RequestStatsServiceClient
	httpClient     *http.Client
	trackURL       string
	timeout        time.Duration
	specURLs       map[string]struct{}
	ignorePrefixes []string
}

type trackHTTPBody struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	Service        string `json:"service"`
	StatusCode     int32  `json:"statusCode"`
	DurationMs     int64  `json:"durationMs"`
	ProfileID      uint64 `json:"profileId"`
	OrganizationID uint64 `json:"organizationId"`
	IsSpecAudit    bool   `json:"isSpecAudit"`
}

// New tạo tracker. trackURL (nếu set) ghi vào v4 request_audit/stats; ngược lại dùng hub conn.
// ignorePrefixes bổ sung (từ config) được merge với defaultIgnorePrefixes.
func New(conn *grpc.ClientConn, specURLs, ignorePrefixes []string, trackURL string) (*Tracker, error) {
	trackURL = strings.TrimSpace(trackURL)
	if trackURL == "" && conn == nil {
		return nil, fmt.Errorf("hub connection or track_url is required for request stats tracker")
	}
	set := make(map[string]struct{}, len(specURLs))
	for _, raw := range specURLs {
		path := strings.TrimSpace(raw)
		if path == "" {
			continue
		}
		set[path] = struct{}{}
	}
	t := &Tracker{
		timeout:        2 * time.Second,
		specURLs:       set,
		ignorePrefixes: mergeIgnorePrefixes(ignorePrefixes),
		trackURL:       trackURL,
		httpClient:     &http.Client{Timeout: 2 * time.Second},
	}
	if trackURL == "" {
		t.client = hubpb.NewRequestStatsServiceClient(conn)
	}
	return t, nil
}

// Middleware sau mỗi request tạo goroutine bắn sang v4/hub để đếm.
func Middleware(tracker *Tracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		if tracker == nil || (tracker.trackURL == "" && tracker.client == nil) {
			return
		}

		path := c.Request.URL.Path
		if tracker.shouldSkip(path) {
			return
		}

		method := c.Request.Method
		statusCode := int32(c.Writer.Status())
		durationMs := time.Since(start).Milliseconds()
		normalizedPath := normalizePath(path)
		serviceName := extractService(normalizedPath)
		isSpecAudit := tracker.isSpecURL(path) || tracker.isSpecURL(normalizedPath)

		var profileID uint64
		var organizationID uint64
		if principal, ok := _httpauth.PrincipalFromContext(c.Request.Context()); ok && principal.ProfileID > 0 {
			profileID = principal.ProfileID
			if principal.OrganizationID != nil {
				organizationID = *principal.OrganizationID
			}
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), tracker.timeout)
			defer cancel()
			if err := tracker.track(ctx, trackHTTPBody{
				Method:         method,
				Path:           normalizedPath,
				Service:        serviceName,
				StatusCode:     statusCode,
				DurationMs:     durationMs,
				ProfileID:      profileID,
				OrganizationID: organizationID,
				IsSpecAudit:    isSpecAudit,
			}); err != nil {
				log.Printf("request_stats track error: %v", err)
			}
		}()
	}
}

func (t *Tracker) track(ctx context.Context, body trackHTTPBody) error {
	if t == nil {
		return nil
	}
	if t.trackURL != "" {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.trackURL, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := t.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("track_url status %d", resp.StatusCode)
		}
		return nil
	}
	if t.client == nil {
		return fmt.Errorf("request stats client is nil")
	}
	_, err := t.client.TrackRequest(ctx, &hubpb.TrackRequestRequest{
		Method:         body.Method,
		Path:           body.Path,
		Service:        body.Service,
		StatusCode:     body.StatusCode,
		DurationMs:     body.DurationMs,
		ProfileId:      body.ProfileID,
		OrganizationId: body.OrganizationID,
		IsSpecAudit:    body.IsSpecAudit,
	})
	return err
}

// defaultIgnorePrefixes: path bắt đầu bằng prefix này sẽ không bị đếm.
// Gồm API dashboard/thống kê hệ thống — tránh tự đếm khi admin mở màn audit.
var defaultIgnorePrefixes = []string{
	"/swagger",
	"/swagger-doc",
	"/healthz",
	"/readyz",
	"/v2/hub/request-stats",        // /dashboard, /daily, /users/*, /audit, /track
	"/v2/hub/keyword-search-stats", // thống kê từ khóa
	"/v2/hub/search-latency",       // /dashboard, /daily
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

// normalizePath thay segment số/uuid/placeholder lỗi (undefined, null, …) bằng :id
// để gom cardinality — kể cả khi client gửi id lỗi vẫn đếm đúng một path pattern.
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

// extractService lấy service từ /v2/{service}/... hoặc /v1/{service}/...
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

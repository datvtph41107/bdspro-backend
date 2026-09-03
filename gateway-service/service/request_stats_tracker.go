package service

import (
	"common"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	_jwt "common/jwt"
	hubpb "pb/types/hub"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	numericPathSegment = regexp.MustCompile(`^[0-9]+$`)
	uuidPathSegment    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	specURLsOnce sync.Once
	specURLSet   map[string]struct{}
)

// RequestStatsTracker gửi thống kê request sang hub (fire-and-forget).
type RequestStatsTracker struct {
	conn    *grpc.ClientConn
	client  hubpb.RequestStatsServiceClient
	timeout time.Duration
}

// NewRequestStatsTracker tạo client gRPC tới hub request-stats.
func NewRequestStatsTracker(env string, opts []grpc.DialOption) (*RequestStatsTracker, error) {
	port := viper.GetString("service.port.hub")
	endpoint := fmt.Sprintf("%s:%s", common.GetPrefixProtobufUrl(env, "hub"), port)
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial hub request stats service: %w", err)
	}
	return &RequestStatsTracker{
		conn:    conn,
		client:  hubpb.NewRequestStatsServiceClient(conn),
		timeout: 2 * time.Second,
	}, nil
}

// Close giải phóng kết nối hub.
func (t *RequestStatsTracker) Close() error {
	if t == nil || t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

// RequestStatsMiddleware sau mỗi request tạo goroutine bắn sang hub để đếm.
func RequestStatsMiddleware(tracker *RequestStatsTracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		if tracker == nil || tracker.client == nil {
			return
		}

		path := c.Request.URL.Path
		if shouldSkipRequestStats(path) {
			return
		}

		method := c.Request.Method
		statusCode := int32(c.Writer.Status())
		durationMs := time.Since(start).Milliseconds()
		normalizedPath := normalizeRequestPath(path)
		serviceName := extractServiceFromPath(normalizedPath)
		isSpecAudit := isSpecURL(path) || isSpecURL(normalizedPath)

		var profileID uint64
		var organizationID uint64
		if claims := _jwt.GetPrincipal(c); claims != nil && claims.ProfileId > 0 {
			profileID = claims.ProfileId
			if claims.OrganizationId != nil {
				organizationID = *claims.OrganizationId
			}
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), tracker.timeout)
			defer cancel()
			_, err := tracker.client.TrackRequest(ctx, &hubpb.TrackRequestRequest{
				Method:         method,
				Path:           normalizedPath,
				Service:        serviceName,
				StatusCode:     statusCode,
				DurationMs:     durationMs,
				ProfileId:      profileID,
				OrganizationId: organizationID,
				IsSpecAudit:    isSpecAudit,
			})
			if err != nil {
				log.Printf("request_stats track error: %v", err)
			}
		}()
	}
}

// requestStatsIgnorePrefixes: path bắt đầu bằng prefix này sẽ không bị đếm.
var requestStatsIgnorePrefixes = []string{
	"/swagger",
	"/swagger-doc",
	"/v2/hub/request-stats",        // không đếm chính API thống kê request
	"/v2/hub/keyword-search-stats", // không đếm API thống kê từ khóa
	"/v2/hub/search-latency",       // không đếm API thống kê latency search
}

func shouldSkipRequestStats(path string) bool {
	if path == "" {
		return true
	}
	for _, prefix := range requestStatsIgnorePrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func loadSpecURLs() map[string]struct{} {
	specURLsOnce.Do(func() {
		specURLSet = make(map[string]struct{})
		for _, raw := range viper.GetStringSlice("request_stats.spec_urls") {
			path := strings.TrimSpace(raw)
			if path == "" {
				continue
			}
			specURLSet[path] = struct{}{}
		}
	})
	return specURLSet
}

// isSpecURL kiểm tra path có nằm trong request_stats.spec_urls không.
func isSpecURL(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	_, ok := loadSpecURLs()[path]
	return ok
}

// normalizeRequestPath thay segment số/uuid bằng :id để gom cardinality.
func normalizeRequestPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/"
	}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "" {
			continue
		}
		if numericPathSegment.MatchString(part) || uuidPathSegment.MatchString(part) {
			parts[i] = ":id"
		}
	}
	normalized := strings.Join(parts, "/")
	if normalized == "" {
		return "/"
	}
	return normalized
}

// extractServiceFromPath lấy service từ /v2/{service}/... hoặc /v1/{service}/...
func extractServiceFromPath(path string) string {
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

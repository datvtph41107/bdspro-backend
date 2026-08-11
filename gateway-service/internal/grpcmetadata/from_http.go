package grpcmetadata

import (
	"common/identitygrpc"
	"common/requestgrpc"
	"context"
	"net"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc/metadata"
)

const createGeneratedReportPath = "/v2/tqd/map-workspace/reports"

type HTTPMetadataAdapter func(context.Context, *http.Request) metadata.MD

// DuringReportMigration keeps every non-target request on the legacy metadata
// bridge. The target bridge is intentionally narrow because it omits legacy
// plan/quota fields and raw credentials.
func DuringReportMigration(
	targetEnabled bool,
	legacy HTTPMetadataAdapter,
) HTTPMetadataAdapter {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		if targetEnabled && isCreateGeneratedReportRequest(r) {
			return FromHTTPRequest(ctx, r)
		}
		if legacy != nil {
			return legacy(ctx, r)
		}
		return metadata.MD{}
	}
}

func ReportTargetEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_REPORT_TARGET_ENABLED")))
	return value == "1" || value == "true" || value == "yes"
}

func isCreateGeneratedReportRequest(r *http.Request) bool {
	if r == nil || r.URL == nil || r.Method != http.MethodPost {
		return false
	}
	return strings.TrimSuffix(r.URL.Path, "/") == createGeneratedReportPath
}

/**
 * FromHTTPRequest chuyển canonical request context sang gRPC metadata.
 *
 * Gateway không đọc lại JWT để dựng Actor tại đây. Actor/Caller phải được
 * auth middleware bind trước đó.
 */
func FromHTTPRequest(_ context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}
	if r == nil {
		return md
	}

	ctx := r.Context()
	md = requestgrpc.AppendFromContext(ctx, md)
	md = identitygrpc.AppendFromContext(ctx, md)

	appendText(md, "client-ip", clientIP(r))
	appendText(md, "user-agent", r.UserAgent())
	appendText(md, "device-id", r.Header.Get("device-id"))
	appendText(md, "if-none-match", r.Header.Get("If-None-Match"))
	appendText(md, "x-qhpro-telemetry-key", r.Header.Get("X-QHPro-Telemetry-Key"))

	return md
}

func appendText(md metadata.MD, key, value string) {
	value = strings.TrimSpace(value)
	if value != "" {
		md.Set(key, value)
	}
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first := strings.TrimSpace(strings.Split(forwarded, ",")[0]); first != "" {
			return first
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

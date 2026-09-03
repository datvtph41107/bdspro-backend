package grpcmetadata

import (
	_rpc "common/rpc"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

// FromHTTPRequest serializes the Gateway-owned canonical context into the
// outgoing gRPC metadata envelope. The canonical RPC client interceptor later
// replaces protected identity fields and signs the envelope.
func FromHTTPRequest(_ context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}
	if r == nil {
		return md
	}
	ctx := r.Context()
	md = _rpc.AppendRequestFromContext(ctx, md)
	md = _rpc.AppendIdentityFromContext(ctx, md)

	appendText(md, "client-ip", clientIP(r))
	appendText(md, "user-agent", r.UserAgent())
	appendText(md, "device-id", r.Header.Get("device-id"))
	appendText(md, "if-none-match", r.Header.Get("If-None-Match"))
	appendText(md, "x-qhpro-telemetry-key", r.Header.Get("X-QHPro-Telemetry-Key"))
	appendText(md, "method", r.Method)
	if r.URL != nil && r.URL.Path == "/v2/payment/sepay/webhook" {
		raw := strings.TrimSpace(r.Header.Get("Authorization"))
		if raw != "" {
			sum := sha256.Sum256([]byte(raw))
			md.Set(_rpc.ProviderCredentialDigestMetadataKey, hex.EncodeToString(sum[:]))
		}
	}
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

// Package httperror owns translation from downstream RPC failures to Gateway
// HTTP responses. It does not authenticate callers or decide service business state.
package httperror

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	sharepb "pb/types/shared"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	dialTCPPortRe           = regexp.MustCompile(`(?i)dial tcp (?:\[[^\]]+\]|[^:\s]+):(\d{2,5})`)
	connectionRefusedPortRe = regexp.MustCompile(`:(\d{2,5}): connect: connection refused`)
)

// WriteGRPC translates one downstream gRPC error into the Gateway HTTP contract.
func WriteGRPC(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("Gateway downstream gRPC error", err)

	grpcErr, ok := status.FromError(err)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	for _, detail := range grpcErr.Details() {
		if info, ok := detail.(*sharepb.ErrorResponse); ok {
			fmt.Printf("Custom code = %d, message = %s\n", info.Code, info.Message)
			// Existing client contract: business errors are returned with HTTP 200
			// and clients inspect the response body code.
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    info.Code,
				"message": info.Message,
				"second":  info.Second,
			})
			return
		}
	}

	w.WriteHeader(statusFromCode(grpcErr.Code()))
	response := map[string]interface{}{
		"code":    grpcErr.Code(),
		"message": publicMessage(grpcErr.Code(), grpcErr.Message()),
	}
	for _, detail := range grpcErr.Details() {
		switch value := detail.(type) {
		case *errdetails.BadRequest:
			violations := make([]map[string]string, 0, len(value.FieldViolations))
			for _, violation := range value.FieldViolations {
				violations = append(violations, map[string]string{
					"field": violation.Field,
					"error": violation.Description,
				})
			}
			response["errors"] = violations
		case *errdetails.ErrorInfo:
			// Preserve owner-provided semantic error data so HTTP clients can
			// present business states without parsing a human message. Gateway
			// translates transport only and does not interpret the metadata.
			response["reason"] = value.Reason
			response["domain"] = value.Domain
			response["metadata"] = value.Metadata
		}
	}
	_ = json.NewEncoder(w).Encode(response)
}

// publicMessage preserves actionable business failures while preventing edge
// responses from exposing internal service addresses, transport details, or
// implementation failures. The complete error remains in the Gateway log with
// the request identity added by middleware.
func publicMessage(code codes.Code, message string) string {
	if friendly, ok := friendlyUpstreamBusyMessage(message); ok {
		return friendly
	}
	switch code {
	case codes.Unavailable:
		// Owner-provided Unavailable messages stay visible. Only dial/transport
		// failures are rewritten to the busy contract above.
		return message
	case codes.DeadlineExceeded:
		return "upstream request timed out"
	case codes.Internal, codes.Unknown, codes.DataLoss:
		return "internal server error"
	default:
		return message
	}
}

func isUpstreamConnectionError(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "error while dialing") ||
		strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "dial tcp") ||
		strings.Contains(lower, "all subconns are in transientfailure") ||
		(strings.Contains(lower, "connection error") && strings.Contains(lower, "transport"))
}

func extractUpstreamPort(msg string) string {
	if m := dialTCPPortRe.FindStringSubmatch(msg); len(m) == 2 {
		return m[1]
	}
	if m := connectionRefusedPortRe.FindStringSubmatch(msg); len(m) == 2 {
		return m[1]
	}
	return ""
}

func friendlyUpstreamBusyMessage(msg string) (string, bool) {
	if !isUpstreamConnectionError(msg) {
		return msg, false
	}
	if port := extractUpstreamPort(msg); port != "" {
		return fmt.Sprintf("máy chủ đang bận vui lòng thử lại sau (%s)", port), true
	}
	return "máy chủ đang bận vui lòng thử lại sau", true
}

func statusFromCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Aborted:
		return http.StatusConflict
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// HandleGRPC adapts WriteGRPC to grpc-gateway's error-handler contract.
func HandleGRPC(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	WriteGRPC(w, err)
}

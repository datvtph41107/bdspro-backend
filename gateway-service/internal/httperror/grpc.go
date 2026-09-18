// Package httperror owns translation from downstream RPC failures to Gateway
// HTTP responses. It does not authenticate callers or decide service business state.
package httperror

import (
	_errors "common/errors"
	"common/request"
	"common/requestlog"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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
	WriteGRPCContext(context.Background(), w, err)
}

// WriteGRPCContext chỉ đổi giao thức; không đoán nghiệp vụ từ message.
// code giữ namespace gRPC 1..16; legacy_code được giữ riêng để không trộn hai bộ mã.
func WriteGRPCContext(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	grpcErr, ok := status.FromError(err)
	if !ok || grpcErr.Code() == codes.OK {
		grpcErr = status.New(codes.Internal, "internal server error")
	}
	code := grpcErr.Code()
	message := grpcErr.Message()
	response := map[string]interface{}{}
	for _, detail := range grpcErr.Details() {
		if info, ok := detail.(*sharepb.ErrorResponse); ok {
			// Lỗi HTTP cũ được nâng đúng transport status; không còn trả 200 khi lỗi.
			code = _errors.RPCCodeForHTTP(info.Code)
			response["legacy_code"] = info.Code
			response["second"] = info.Second
			if info.Message != "" {
				message = info.Message
			}
		}
	}
	httpStatus := statusFromCode(code)
	response["code"] = code
	response["message"] = publicMessage(code, message)
	response["status"] = httpStatus
	response["reason"] = code.String()
	requestID, _ := request.RequestIDFromContext(ctx)
	operationID, _ := request.OperationIDFromContext(ctx)
	if requestID == "" {
		requestID = w.Header().Get("X-Request-Id")
	}
	if operationID == "" {
		operationID = w.Header().Get("X-Operation-Id")
	}
	if requestID != "" {
		response["request_id"] = requestID
	}
	if operationID != "" {
		response["operation_id"] = operationID
	}
	// Chỉ business 4xx được mang detail công khai; không lộ metadata hạ tầng ở 5xx.
	if httpStatus < 500 {
		for _, detail := range grpcErr.Details() {
			switch value := detail.(type) {
			case *errdetails.BadRequest:
				violations := make([]map[string]string, 0, len(value.FieldViolations))
				for _, violation := range value.FieldViolations {
					violations = append(violations, map[string]string{"field": violation.Field, "error": violation.Description})
				}
				response["errors"] = violations
			case *errdetails.ErrorInfo:
				response["reason"] = value.Reason
				response["domain"] = value.Domain
				response["metadata"] = value.Metadata
			}
		}
	}
	level := slog.LevelWarn
	if httpStatus >= 500 {
		level = slog.LevelError
	}
	route, _ := runtime.HTTPPathPattern(ctx)
	method, _ := runtime.RPCMethod(ctx)
	requestlog.FromContext(ctx, nil).Log(ctx, level, "Downstream request rejected", "http_route", route, "grpc_method", method, "category", "dependency", "event_name", "gateway.downstream.failed", "grpc_code", code.String(), "http_status_code", httpStatus, "error_reason", response["reason"], "error_domain", response["domain"])
	w.WriteHeader(httpStatus)
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
		if isInfrastructureError(message) || strings.TrimSpace(message) == "" {
			return "internal server error"
		}
		// Owner-provided business/diagnostic messages (via ReturnError) stay visible.
		return message
	default:
		return message
	}
}

func isInfrastructureError(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "sqlstate") ||
		strings.Contains(lower, "sql:") ||
		strings.Contains(lower, "sql ") ||
		strings.Contains(lower, "pq:") ||
		strings.Contains(lower, "dial tcp") ||
		strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "runtime error") ||
		strings.Contains(lower, "column \"") ||
		strings.Contains(lower, "relation \"") ||
		strings.Contains(lower, "dsn") ||
		strings.Contains(lower, "password=") ||
		strings.Contains(lower, "stack")
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

func statusFromCode(code codes.Code) int { return _errors.HTTPStatusForRPC(code) }

// HandleGRPC giữ context đã bind từ HTTP để truy vết request xuyên service.
func HandleGRPC(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	WriteGRPCContext(ctx, w, err)
}


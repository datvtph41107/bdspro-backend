// Package httperror owns translation from downstream RPC failures to Gateway
// HTTP responses. It does not authenticate callers or decide service business state.
package httperror

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	_logging "common/logging"
	_httpresponse "gateway/internal/httpresponse"
	sharepb "pb/types/shared"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WriteGRPC is retained for tests and compatibility callers that do not carry a
// request context. Production grpc-gateway handling uses WriteGRPCContext.
func WriteGRPC(w http.ResponseWriter, err error) {
	WriteGRPCContext(context.Background(), w, err)
}

// WriteGRPCContext translates one downstream gRPC error into the canonical
// public HTTP problem contract and emits one structured boundary event.
func WriteGRPCContext(ctx context.Context, w http.ResponseWriter, err error) {
	grpcErr, ok := status.FromError(err)
	if !ok {
		_logging.WithChannel(ctx, "http").With(
			slog.String("component", "gateway.httperror"),
		).ErrorContext(ctx, "downstream non-gRPC error",
			slog.String("event_name", "gateway.downstream.error"),
			slog.Any("error", err),
			slog.Int("http_status_code", http.StatusInternalServerError),
		)
		_httpresponse.WriteProblem(ctx, w, _httpresponse.Problem{
			Type:   "urn:qhpro:error:rpc.internal",
			Title:  http.StatusText(http.StatusInternalServerError),
			Status: http.StatusInternalServerError,
			Detail: "internal server error",
			Code:   "rpc.internal",
		})
		return
	}

	// Migration-only compatibility adapter. Existing legacy protobuf business
	// errors historically return HTTP 200 and a numeric body code. New service
	// code must not create sharepb.ErrorResponse; typed common/fault errors use
	// the canonical RFC 9457 path below. Remove this adapter when legacy emitters
	// reach zero.
	for _, detail := range grpcErr.Details() {
		if info, ok := detail.(*sharepb.ErrorResponse); ok {
			writeLegacyBusinessError(ctx, w, grpcErr, info)
			return
		}
	}

	problem := problemFromGRPC(grpcErr)
	logProblem(ctx, grpcErr, problem)
	_httpresponse.WriteProblem(ctx, w, problem)
}

func writeLegacyBusinessError(ctx context.Context, w http.ResponseWriter, grpcErr *status.Status, info *sharepb.ErrorResponse) {
	_logging.WithChannel(ctx, "http").With(
		slog.String("component", "gateway.httperror"),
	).WarnContext(ctx, "legacy downstream business error",
		slog.String("event_name", "gateway.downstream.legacy_business_error"),
		slog.String("grpc_code", grpcErr.Code().String()),
		slog.Int64("legacy_error_code", int64(info.Code)),
		slog.Bool("legacy_http_200_contract", true),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    info.Code,
		"message": info.Message,
		"second":  info.Second,
	})
}

func problemFromGRPC(grpcErr *status.Status) _httpresponse.Problem {
	httpStatus := statusFromCode(grpcErr.Code())
	problem := _httpresponse.Problem{
		Status: httpStatus,
		Title:  http.StatusText(httpStatus),
		Detail: publicGRPCDetail(grpcErr, httpStatus),
		Code:   fallbackProblemCode(grpcErr.Code()),
	}

	for _, detail := range grpcErr.Details() {
		switch value := detail.(type) {
		case *errdetails.BadRequest:
			problem.Errors = make([]_httpresponse.FieldProblem, 0, len(value.FieldViolations))
			for _, violation := range value.FieldViolations {
				problem.Errors = append(problem.Errors, _httpresponse.FieldProblem{
					Field:  violation.Field,
					Detail: violation.Description,
				})
			}
		case *errdetails.ErrorInfo:
			problem.Reason = value.Reason
			problem.Domain = value.Domain
			if len(value.Metadata) > 0 {
				problem.Metadata = make(map[string]string, len(value.Metadata))
				for key, metadataValue := range value.Metadata {
					problem.Metadata[key] = metadataValue
				}
				if errorCode := strings.TrimSpace(value.Metadata["error_code"]); errorCode != "" {
					problem.Code = errorCode
				}
			}
			if problem.Code == fallbackProblemCode(grpcErr.Code()) && strings.TrimSpace(value.Reason) != "" {
				problem.Code = codeFromReason(value.Reason)
			}
		}
	}

	problem.Type = problemType(problem.Code)
	return problem
}

func logProblem(ctx context.Context, grpcErr *status.Status, problem _httpresponse.Problem) {
	level := slog.LevelWarn
	if problem.Status >= http.StatusInternalServerError {
		level = slog.LevelError
	}
	logger := _logging.WithChannel(ctx, "http").With(
		slog.String("component", "gateway.httperror"),
	)
	logger.LogAttrs(ctx, level, "downstream gRPC request failed",
		slog.String("event_name", "gateway.downstream.grpc_error"),
		slog.String("grpc_code", grpcErr.Code().String()),
		slog.Int("http_status_code", problem.Status),
		slog.String("error_code", problem.Code),
		slog.String("error_reason", problem.Reason),
	)
}

func publicGRPCDetail(grpcErr *status.Status, httpStatus int) string {
	if httpStatus < http.StatusInternalServerError {
		return grpcErr.Message()
	}
	switch grpcErr.Code() {
	case codes.Unavailable:
		return "service unavailable"
	case codes.DeadlineExceeded:
		return "request deadline exceeded"
	default:
		return "internal server error"
	}
}

func problemType(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return "about:blank"
	}
	var result strings.Builder
	for _, value := range code {
		switch {
		case value >= 'a' && value <= 'z', value >= 'A' && value <= 'Z', value >= '0' && value <= '9', value == '.', value == '-', value == '_', value == ':':
			result.WriteRune(value)
		default:
			result.WriteByte('-')
		}
	}
	return "urn:qhpro:error:" + result.String()
}

func codeFromReason(reason string) string {
	reason = strings.TrimSpace(strings.ToLower(reason))
	reason = strings.ReplaceAll(reason, "_", ".")
	return reason
}

func fallbackProblemCode(code codes.Code) string {
	switch code {
	case codes.InvalidArgument:
		return "rpc.invalid_argument"
	case codes.Unauthenticated:
		return "rpc.unauthenticated"
	case codes.PermissionDenied:
		return "rpc.permission_denied"
	case codes.NotFound:
		return "rpc.not_found"
	case codes.AlreadyExists:
		return "rpc.already_exists"
	case codes.Aborted:
		return "rpc.aborted"
	case codes.FailedPrecondition:
		return "rpc.failed_precondition"
	case codes.ResourceExhausted:
		return "rpc.resource_exhausted"
	case codes.Canceled:
		return "rpc.canceled"
	case codes.DeadlineExceeded:
		return "rpc.deadline_exceeded"
	case codes.Unimplemented:
		return "rpc.unimplemented"
	case codes.Unavailable:
		return "rpc.unavailable"
	default:
		return "rpc.internal"
	}
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

// HandleGRPC adapts WriteGRPCContext to grpc-gateway's error-handler contract.
func HandleGRPC(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	WriteGRPCContext(ctx, w, err)
}

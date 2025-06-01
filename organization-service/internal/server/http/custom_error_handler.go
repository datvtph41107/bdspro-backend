package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func extractErrorDetails(st *status.Status) []map[string]interface{} {
	details := make([]map[string]interface{}, 0)

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.GetFieldViolations() {
				details = append(details, map[string]interface{}{
					"field":       violation.GetField(),
					"description": violation.GetDescription(),
					"type":        "FIELD_VIOLATION",
				})
			}
		case *errdetails.RetryInfo:
			details = append(details, map[string]interface{}{
				"type":       "RETRY_INFO",
				"retryDelay": t.GetRetryDelay().String(),
			})
		case *errdetails.DebugInfo:
			details = append(details, map[string]interface{}{
				"type":         "DEBUG_INFO",
				"stackEntries": t.GetStackEntries(),
				"detail":       t.GetDetail(),
			})

		}
	}
	return details
}

func CustomErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s := status.Convert(err)

	w.Header().Set("Content-Type", "application/json")

	status := HTTPStatusFromCode(s.Code())
	w.WriteHeader(status)

	errorResponse := map[string]interface{}{
		"code":    s.Code(),
		"message": s.Message(),
		"details": extractErrorDetails(s),
	}

	json.NewEncoder(w).Encode(errorResponse)
}


func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case custom_error.InvalidRequestCode:
		return http.StatusBadRequest
	default:
		grpclog.Warningf("Unknown gRPC error code: %v", code)
		return http.StatusInternalServerError
	}
}
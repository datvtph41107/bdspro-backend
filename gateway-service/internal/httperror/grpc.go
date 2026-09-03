// Package httperror owns translation from downstream RPC failures to Gateway
// HTTP responses. It does not authenticate callers or decide service business state.
package httperror

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	sharepb "pb/types/shared"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		"message": grpcErr.Message(),
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

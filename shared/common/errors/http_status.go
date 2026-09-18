package _errors

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

// RPCCodeForHTTP là bảng chuyển status tường minh cho API lỗi cũ dùng số HTTP.
// Mã nghiệp vụ riêng không được cast thành gRPC code vì gRPC chỉ có 0..16.
func RPCCodeForHTTP(code int32) codes.Code {
	// Auth/OTP business codes (1000-1999) are client-facing validation failures.
	if code >= 1000 && code < 2000 {
		return codes.InvalidArgument
	}
	switch code {
	case 400, 422:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 408:
		return codes.Canceled
	case 409:
		return codes.AlreadyExists
	case 412:
		return codes.FailedPrecondition
	case 429:
		return codes.ResourceExhausted
	case 501:
		return codes.Unimplemented
	case 502, 503:
		return codes.Unavailable
	case 504:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}
func HTTPStatusForRPC(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
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

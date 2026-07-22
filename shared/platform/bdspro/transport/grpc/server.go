package bdspro_grpc

import (
	apperr "bdspro/core/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCStatus converts an AppError to a gRPC status.
// It maps Kind to gRPC codes and attaches custom details if available.
func ToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}
	e, ok := apperr.As(err)
	if !ok {
		// Unknown error: internal
		return status.Error(codes.Internal, "internal server error")
	}

	var code codes.Code
	switch e.Kind() {
	case apperr.KindValidation:
		code = codes.InvalidArgument
	case apperr.KindBusiness:
		code = mapBusinessCode(e.Code())
	case apperr.KindInternal:
		code = codes.Internal
	case apperr.KindExternal:
		code = codes.Unavailable
	default:
		code = codes.Unknown
	}

	// Build status
	st := status.New(code, e.Message())

	// Add details (e.g., BadRequest for validation, custom AppError detail)
	// For now, return simple status. Can extend later.
	return st.Err()
}

// FromGRPCStatus converts a gRPC status error to an AppError.
func FromGRPCStatus(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return apperr.Internal("grpc.from_status", err)
	}

	// Try to extract AppError from details (if we added custom detail)
	// For now, map based on gRPC code.
	switch st.Code() {
	case codes.InvalidArgument:
		return apperr.ErrInvalidArgument
	case codes.NotFound:
		return apperr.ErrNotFound
	case codes.AlreadyExists:
		return apperr.ErrAlreadyExists
	case codes.PermissionDenied:
		return apperr.ErrForbidden
	case codes.Unauthenticated:
		return apperr.ErrUnauthorized
	case codes.Unavailable:
		return apperr.ErrUnavailable
	case codes.DeadlineExceeded:
		return apperr.ErrTimeout
	default:
		// fallback to internal
		return apperr.Internal("grpc.from_status", err)
	}
}

func mapBusinessCode(code apperr.Code) codes.Code {
	switch code {
	case apperr.CodeNotFound:
		return codes.NotFound
	case apperr.CodeAlreadyExists:
		return codes.AlreadyExists
	case apperr.CodeForbidden:
		return codes.PermissionDenied
	case apperr.CodeUnauthorized:
		return codes.Unauthenticated
	case apperr.CodeConflict:
		return codes.Aborted
	default:
		return codes.InvalidArgument
	}
}

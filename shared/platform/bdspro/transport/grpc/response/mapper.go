package grpc_resp

import (
	apperr "bdspro/core/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCStatus converts AppError to gRPC status
func ToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}
	e, ok := apperr.As(err)
	if !ok {
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
	return status.New(code, e.Message()).Err()
}

// FromGRPCStatus converts gRPC status to AppError
func FromGRPCStatus(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return apperr.Internal("grpc.from_status", err)
	}

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

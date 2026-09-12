package handler

import (
	"errors"

	_routes "common/routes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapRefreshTokenError preserves the authentication boundary at the gRPC
// transport. An expired or stale session is expected client state and must
// never be presented as an internal server failure.
func mapRefreshTokenError(err error) error {
	if err == nil {
		return nil
	}

	var applicationError *_routes.Except
	if errors.As(err, &applicationError) {
		code := codes.Internal
		switch applicationError.Code {
		case 400:
			code = codes.InvalidArgument
		case 401:
			code = codes.Unauthenticated
		case 403:
			code = codes.PermissionDenied
		case 404:
			code = codes.NotFound
		case 409:
			code = codes.AlreadyExists
		case 429:
			code = codes.ResourceExhausted
		case 503:
			code = codes.Unavailable
		}
		return status.Error(code, applicationError.Message)
	}

	if grpcStatus, ok := status.FromError(err); ok && grpcStatus.Code() != codes.Unknown {
		return err
	}

	return status.Error(codes.Internal, "không thể làm mới phiên đăng nhập")
}

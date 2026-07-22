package grpc_interceptor

import (
	apperr "bdspro/core/pkg/errors"
	grpc_resp "bdspro/transport/grpc/response"
	"context"
	"fmt"
	"runtime/debug"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = apperr.Internal("grpc.panic", fmt.Errorf("%v", r)).
					WithField("stack", "panic_stack", string(debug.Stack()), nil)
				err = grpc_resp.ToGRPCStatus(err)
			}
		}()

		resp, err = handler(ctx, req)
		if err != nil {
			return nil, grpc_resp.ToGRPCStatus(err)
		}
		return resp, nil
	}
}

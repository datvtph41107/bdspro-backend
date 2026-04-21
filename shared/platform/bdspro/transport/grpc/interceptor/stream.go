package grpc_interceptor

import (
	"fmt"
	"runtime/debug"

	apperr "bdspro/core/pkg/errors"
	grpc_resp "bdspro/transport/grpc/response"

	"google.golang.org/grpc"
)

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = apperr.Internal("grpc.panic", fmt.Errorf("%v", r)).
					WithField("stack", "panic_stack", string(debug.Stack()), nil)
				err = grpc_resp.ToGRPCStatus(err)
			}
		}()

		err = handler(srv, ss)
		if err != nil {
			return grpc_resp.ToGRPCStatus(err)
		}
		return nil
	}
}

package interceptors

import (
	"context"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
)

func UnaryLoggerInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		logger.Infof("Unary gRPC method: %s, request: %v, duration: %v, error: %v", info.FullMethod, req, time.Since(start), err)

		return resp, err
	}
}

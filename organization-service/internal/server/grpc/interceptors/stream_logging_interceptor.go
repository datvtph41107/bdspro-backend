package interceptors

import (
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
)

func StreamLoggerInterceptor(logger *flogging.FabricLogger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, ss)

		logger.Infof("Stream gRPC method: %s, isClientStream: %v, isServerStream: %v, duration: %v, error: %v", info.FullMethod, info.IsClientStream, info.IsServerStream, time.Since(start), err)

		return err
	}
}

package interceptors

import (
	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
)

func StreamRecoveryInterceptor(logger *flogging.FabricLogger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("panic recovered in stream interceptor method: %s, panic: %v", info.FullMethod, r)
			}
		}()
		return handler(srv, ss)
	}
}

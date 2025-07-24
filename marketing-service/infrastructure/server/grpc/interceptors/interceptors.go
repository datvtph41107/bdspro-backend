package interceptors

import (
	"context"

	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryLoggerInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infof("gRPC method: %s", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Errorf("gRPC method %s failed: %v", info.FullMethod, err)
		}
		return resp, err
	}
}

func UnaryRecoveryInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Panic recovered in %s: %v", info.FullMethod, r)
				err = status.Errorf(codes.Internal, "Internal server error: %v", r)
			}
		}()
		return handler(ctx, req)
	}
}

func UnaryAuthInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Implement authentication logic
		logger.Infof("Auth interceptor for method: %s", info.FullMethod)
		return handler(ctx, req)
	}
} 
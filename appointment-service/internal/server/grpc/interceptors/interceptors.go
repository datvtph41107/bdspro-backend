package interceptors

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
)

func UnaryLoggerInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infof("Request - Method:%s\n", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Errorf("Error - Method:%s\tError:%v\n", info.FullMethod, err)
		}
		return resp, err
	}
}

func UnaryAuthInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Implement authentication logic
		return handler(ctx, req)
	}
}

func UnaryRecoveryInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
				logger.Errorf("Recovery from panic: %v", r)
			}
		}()
		return handler(ctx, req)
	}
} 
package interceptors

import (
	"context"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/tuandq2112/appointment-service/appconfig"
	utils "github.com/tuandq2112/appointment-service/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func UnaryAuthInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/organizationpb.OrganizationMemberService/SwitchOrganization" {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("grpcgateway-authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokenString := authHeader[0]
		claims, err := utils.DecodeJWTPayload(tokenString)
		if err != nil {
			logger.Error("Invalid token", err)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		newCtx := context.WithValue(ctx, appconfig.USER_CONTEXT, claims)
		return handler(newCtx, req)
	}
}

func UnaryRecoveryInterceptor(logger *flogging.FabricLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("panic recovered in unary interceptor method: %s, panic: %v", info.FullMethod, r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

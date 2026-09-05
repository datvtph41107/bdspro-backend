package interceptors

import (
	"context"
	"payment/internal/requestactor"
	"payment/pkg/utils"
	"slices"

	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(logger *flogging.FabricLogger, whitelist []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if slices.Contains(whitelist, info.FullMethod) {
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
		newCtx := context.WithValue(ctx, requestactor.UserContextKey, claims)
		return handler(newCtx, req)
	}
}

package interceptors

import (
	"common/logging"
	"context"
	"log/slog"
	"payment/internal/requestactor"
	"payment/pkg/utils"
	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(whitelist []string) grpc.UnaryServerInterceptor {
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
			logging.WithComponent(ctx, "grpc.auth").Error(
				"invalid authentication token",
				slog.String("method", info.FullMethod),
				slog.Any("error", err),
			)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		newCtx := context.WithValue(ctx, requestactor.UserContextKey, claims)
		return handler(newCtx, req)
	}
}

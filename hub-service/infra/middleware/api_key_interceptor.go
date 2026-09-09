package middleware

import (
	"context"
	"strings"

	_enum "common/domain/enum"
	_usecase "hub/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewAPIKeyUnaryServerInterceptor builds the protected-method lookup once from
// the process-owned runtime snapshot. Request handling never reads global
// configuration.
func NewAPIKeyUnaryServerInterceptor(
	apiKeyUsecase _usecase.IApiKeyUsecase,
	protectedMethods []string,
) grpc.UnaryServerInterceptor {
	protected := make(map[string]struct{})
	for _, method := range protectedMethods {
		method = strings.TrimSpace(method)
		if method == "" {
			continue
		}
		protected[method] = struct{}{}
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if len(protected) == 0 {
			return handler(ctx, req)
		}
		if _, ok := protected[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		apiKey := extractAPIKey(ctx)
		if apiKey == "" {
			return nil, status.Error(codes.Unauthenticated, "missing api key")
		}
		if apiKeyUsecase == nil {
			return nil, status.Error(codes.Internal, "api key usecase not initialized")
		}

		if _, errDTO := apiKeyUsecase.VerifyApiKey(ctx, apiKey); errDTO != nil {
			switch errDTO.Code {
			case 400:
				return nil, status.Error(codes.InvalidArgument, errDTO.Message)
			case 401:
				return nil, status.Error(codes.Unauthenticated, errDTO.Message)
			case 404:
				return nil, status.Error(codes.PermissionDenied, errDTO.Message)
			default:
				return nil, status.Error(codes.Internal, errDTO.Message)
			}
		}

		ctx = context.WithValue(ctx, _enum.APIKeyKey, apiKey)
		return handler(ctx, req)
	}
}

func extractAPIKey(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(_enum.APIKeyKey).(string); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("x-api-key")
		if len(values) > 0 {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}

package middleware

// ─────────────────────────────────────────────────────────────────────────────
// hub/infra/middleware/apikey_interceptor.go
//
// Middleware này dùng config.AppProperties.Security.ProtectedMethods
// để quyết định method nào cần kiểm tra ApiKey.
//
// ApplinkService KHÔNG có trong ProtectedMethods
// → tự động bypass → không cần thêm logic gì đặc biệt.
// ─────────────────────────────────────────────────────────────────────────────

import (
	"context"
	"strings"

	_enum "common/domain/enum"
	"hub/config"
	_usecase "hub/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewAPIKeyUnaryServerInterceptor tạo interceptor kiểm tra API key.
//
// Logic:
//  1. Đọc danh sách protected methods từ config (whitelist)
//  2. Nếu method KHÔNG có trong protected → bypass (public)
//  3. Nếu method CÓ trong protected → yêu cầu x-api-key header hợp lệ
//
// ApplinkService/CreateApplink và ApplinkService/GetApplinkByCode
// là public — KHÔNG thêm vào config.Security.ProtectedMethods.
func NewAPIKeyUnaryServerInterceptor(apiKeyUsecase _usecase.IApiKeyUsecase) grpc.UnaryServerInterceptor {
	// Build protected methods map một lần khi khởi động
	protected := make(map[string]struct{})
	for _, method := range config.AppProperties.Security.ProtectedMethods {
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
		// Không có protected method nào được config → bypass tất cả
		if len(protected) == 0 {
			return handler(ctx, req)
		}

		// Method không nằm trong whitelist → public, bypass
		if _, ok := protected[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		// ── Method yêu cầu ApiKey ─────────────────────────────────────────
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

// extractAPIKey lấy api key từ context value hoặc gRPC metadata header x-api-key
func extractAPIKey(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Ưu tiên lấy từ context value (đã được set bởi lần trước)
	if val, ok := ctx.Value(_enum.APIKeyKey).(string); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}

	// Fallback: lấy từ gRPC metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("x-api-key")
		if len(values) > 0 {
			return strings.TrimSpace(values[0])
		}
	}

	return ""
}

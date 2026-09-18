package cmd

import (
	"common/operation"
	operationgrpc "common/operation/grpc"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"google.golang.org/protobuf/proto"
	quotausecase "tqd/internal/usecase/quota"
	"tqd/internal/usecase/usage"

	tqdpb "pb/types/tqd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// usesReportTargetRuntime reports whether the canonical Generated Report
// runtime owns this RPC. There is no feature-flag fallback for this method.
func usesReportTargetRuntime(fullMethod string) bool {
	return fullMethod == tqdpb.MapWorkspaceService_CreateGeneratedReport_FullMethodName
}

func usesCommercialQuotaRuntime(fullMethod string) bool {
	return usesReportTargetRuntime(fullMethod) ||
		fullMethod == tqdpb.ParcelService_FindParcelsByPolygon_FullMethodName ||
		fullMethod == tqdpb.ParcelService_GetParcelDetail_FullMethodName
}

func usesActionQuotaRuntime(fullMethod string) bool {
	return fullMethod == tqdpb.ParcelService_FindParcelsByPolygon_FullMethodName ||
		fullMethod == tqdpb.ParcelService_GetParcelDetail_FullMethodName
}

func usesCanonicalTargetRuntime(fullMethod string) bool {
	return usesCommercialQuotaRuntime(fullMethod) ||
		fullMethod == tqdpb.AdminUsageService_ListUsage_FullMethodName ||
		fullMethod == tqdpb.AdminUsageService_ReconcileUsage_FullMethodName ||
		fullMethod == tqdpb.ProfileUsageService_GetMyGeneratedReportUsage_FullMethodName ||
		fullMethod == tqdpb.ProfileUsageService_ListMyUsage_FullMethodName
}

// buildCommercialActionInterceptor accounts successful customer-facing TQD
// reads. Generated Report keeps its stronger transactional acceptance theorem
// and is therefore deliberately excluded here.
func buildCommercialActionInterceptor(service *quotausecase.ActionService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info == nil {
			return nil, status.Error(codes.Internal, "gRPC method info is missing")
		}
		if !usesActionQuotaRuntime(info.FullMethod) {
			return handler(ctx, req)
		}
		if service == nil {
			return nil, status.Error(codes.Unavailable, "commercial quota is unavailable")
		}
		operationCode, ok := operation.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "commercial operation is missing")
		}

		message, ok := req.(proto.Message)
		if !ok || message == nil {
			return nil, status.Error(codes.InvalidArgument, "protobuf request is required")
		}
		payload, err := proto.MarshalOptions{Deterministic: true}.Marshal(message)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid read payload")
		}
		digest := sha256.Sum256(payload)
		var response interface{}
		err = service.Run(ctx, operationCode, hex.EncodeToString(digest[:]), func(actionCtx context.Context) error {
			var actionErr error
			response, actionErr = handler(actionCtx, req)
			return actionErr
		})
		if err != nil {
			return nil, commercialActionStatus(err)
		}
		return response, nil
	}
}

func commercialActionStatus(err error) error {
	var denied *quotausecase.AccessDeniedError
	var exhausted *quotausecase.ExhaustedError
	switch {
	case errors.Is(err, usage.ErrReadCommandConflict):
		return status.Error(codes.AlreadyExists, "Idempotency-Key đã dùng với dữ liệu khác")
	case errors.Is(err, usage.ErrLegacyReadCommand):
		return status.Error(codes.FailedPrecondition, "Thao tác cũ thiếu dấu vân tay dữ liệu; hãy bắt đầu thao tác mới")
	case errors.As(err, &denied), errors.Is(err, quotausecase.ErrAccessDenied):
		return status.Error(codes.PermissionDenied, "gói hiện tại không bao gồm tính năng này")
	case errors.As(err, &exhausted), errors.Is(err, quotausecase.ErrQuotaExceeded):
		return status.Error(codes.ResourceExhausted, "đã sử dụng hết lượt trong chu kỳ hiện tại")
	case errors.Is(err, quotausecase.ErrActionIdentity), errors.Is(err, quotausecase.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "Idempotency-Key hợp lệ là bắt buộc cho thao tác tính quota")
	case errors.Is(err, quotausecase.ErrActionUnavailable):
		return status.Error(codes.Unavailable, "không thể xác minh hoặc ghi nhận quota")
	default:
		// Preserve business statuses returned by the wrapped TQD handler.
		return err
	}
}

// buildReportTargetOperationInterceptor binds the contract-owned operation only
// for the canonical Generated Report RPC. Other TQD RPCs keep their existing
// operation behavior until migrated independently.
func buildReportTargetOperationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil {
			return nil, status.Error(
				codes.Internal,
				"gRPC method info is missing",
			)
		}
		if !usesCommercialQuotaRuntime(info.FullMethod) {
			return handler(ctx, req)
		}

		return operationgrpc.UnaryServerInterceptor(
			ctx,
			req,
			info,
			handler,
		)
	}
}

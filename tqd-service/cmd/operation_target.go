package cmd

import (
	operationgrpc "common/operation/grpc"
	"context"

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

func usesCanonicalTargetRuntime(fullMethod string) bool {
	return usesReportTargetRuntime(fullMethod) ||
		fullMethod == tqdpb.AdminUsageService_ListUsage_FullMethodName ||
		fullMethod == tqdpb.AdminUsageService_ReconcileUsage_FullMethodName ||
		fullMethod == tqdpb.ProfileUsageService_GetMyGeneratedReportUsage_FullMethodName
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
		if !usesReportTargetRuntime(info.FullMethod) {
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

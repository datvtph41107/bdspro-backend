package cmd

import (
	"common/businessoperation"
	"context"
	"tqd/infra/grpcoperation"

	"google.golang.org/grpc"
)

const createGeneratedReportMethod = "/types.tqd.MapWorkspaceService/CreateGeneratedReport"

/**
 * newBusinessOperationInterceptor map gRPC method sang business operation mới.
 */
func newBusinessOperationInterceptor(targetEnabled bool) (grpc.UnaryServerInterceptor, error) {
	if !targetEnabled {
		return func(
			ctx context.Context,
			req interface{},
			_ *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			return handler(ctx, req)
		}, nil
	}

	registry := grpcoperation.NewRegistry()
	if err := registry.Add(
		createGeneratedReportMethod,
		businessoperation.WorkspaceReportGenerate,
	); err != nil {
		return nil, err
	}
	return grpcoperation.UnaryServerInterceptor(registry), nil
}

/**
 * usesTargetQuota trả true khi RPC đã tự xử lý quota trong business flow mới.
 */
func usesTargetQuota(targetEnabled bool, fullMethod string) bool {
	return targetEnabled && fullMethod == createGeneratedReportMethod
}

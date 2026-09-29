package cmd

import (
	"common/operation"
	"context"
	"testing"

	tqdpb "pb/types/tqd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUsesReportTargetRuntimeOwnsGeneratedReportOnly(t *testing.T) {
	t.Parallel()

	method := tqdpb.MapWorkspaceService_CreateGeneratedReport_FullMethodName
	if !usesReportTargetRuntime(method) {
		t.Fatal("canonical Generated Report RPC must always use target runtime")
	}
	if usesReportTargetRuntime(tqdpb.MapWorkspaceService_ListFollowedParcels_FullMethodName) {
		t.Fatal("non-target RPC must stay outside the Report target runtime")
	}
}

func TestReportTargetOperationInterceptorBindsCanonicalSelection(t *testing.T) {
	t.Parallel()

	_, err := buildReportTargetOperationInterceptor()(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: tqdpb.MapWorkspaceService_CreateGeneratedReport_FullMethodName,
		},
		func(ctx context.Context, _ interface{}) (interface{}, error) {
			selected, ok := operation.FromContext(ctx)
			if !ok {
				t.Fatal("canonical operation is missing")
			}
			if selected != operation.Code("workspace.report.generate") {
				t.Fatalf("canonical operation = %q", selected)
			}
			return "ok", nil
		},
	)
	if err != nil {
		t.Fatalf("report target operation interceptor error = %v", err)
	}
}

func TestReportTargetOperationInterceptorLeavesOtherRPCsUntouched(t *testing.T) {
	t.Parallel()

	_, err := buildReportTargetOperationInterceptor()(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: tqdpb.MapWorkspaceService_ListFollowedParcels_FullMethodName,
		},
		func(ctx context.Context, _ interface{}) (interface{}, error) {
			if code, ok := operation.FromContext(ctx); ok {
				t.Fatalf("unexpected canonical operation = %q", code)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReportTargetOperationInterceptorRejectsCanonicalContextConflict(t *testing.T) {
	t.Parallel()

	ctx, err := operation.Bind(
		context.Background(),
		operation.Code("workspace.report.delete"),
	)
	if err != nil {
		t.Fatal(err)
	}

	called := false
	_, err = buildReportTargetOperationInterceptor()(
		ctx,
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: tqdpb.MapWorkspaceService_CreateGeneratedReport_FullMethodName,
		},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)

	if status.Code(err) != codes.Internal {
		t.Fatalf("status code = %v, want Internal; error = %v", status.Code(err), err)
	}
	if called {
		t.Fatal("handler must not run after canonical operation conflict")
	}
}

package cmd

import (
	_operationcatalog "common/operationcatalog"
	_operationcontext "common/operationcontext"

	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestTQDOperationRuntimeSelectsGeneratedReport(t *testing.T) {
	runtime, err := newTQDOperationRuntime()
	if err != nil {
		t.Fatalf("newTQDOperationRuntime() error = %v", err)
	}

	_, err = runtime.UnaryServerInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.MapWorkspaceService/CreateGeneratedReport",
		},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			selection, ok := _operationcontext.FromContext(ctx)
			if !ok {
				t.Fatal("operation is missing")
			}
			if selection.Code != _operationcatalog.WorkspaceReportGenerateOperation {
				t.Fatalf("operation = %q", selection.Code)
			}
			return struct{}{}, nil
		},
	)
	if err != nil {
		t.Fatalf("UnaryServerInterceptor() error = %v", err)
	}
}

func TestTQDOperationRuntimeLeavesUnapprovedRPCUnselected(t *testing.T) {
	runtime, err := newTQDOperationRuntime()
	if err != nil {
		t.Fatalf("newTQDOperationRuntime() error = %v", err)
	}

	_, err = runtime.UnaryServerInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.ParcelService/FindParcelByLocation",
		},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			if selection, ok := _operationcontext.FromContext(ctx); ok {
				t.Fatalf("unexpected operation = %+v", selection)
			}
			return struct{}{}, nil
		},
	)
	if err != nil {
		t.Fatalf("UnaryServerInterceptor() error = %v", err)
	}
}

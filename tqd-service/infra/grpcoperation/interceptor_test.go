package grpcoperation

import (
	"common/businessoperation"
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestUnaryServerInterceptorBindsOperation(t *testing.T) {
	registry := NewRegistry()
	method := "/types.tqd.MapWorkspaceService/CreateGeneratedReport"
	if err := registry.Add(method, businessoperation.WorkspaceReportGenerate); err != nil {
		t.Fatal(err)
	}

	interceptor := UnaryServerInterceptor(registry)
	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: method},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			operation, ok := businessoperation.FromContext(ctx)
			if !ok || operation != businessoperation.WorkspaceReportGenerate {
				t.Fatalf("operation = %q, %v", operation, ok)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor error = %v", err)
	}
}

func TestUnaryServerInterceptorLeavesUnknownMethodUnselected(t *testing.T) {
	interceptor := UnaryServerInterceptor(NewRegistry())
	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/types.tqd.ParcelService/GetParcel"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			if _, ok := businessoperation.FromContext(ctx); ok {
				t.Fatal("unknown method should not bind an operation")
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor error = %v", err)
	}
}

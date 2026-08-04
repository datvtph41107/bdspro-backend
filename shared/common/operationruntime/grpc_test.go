package operationruntime

import (
	_operationbinding "common/operationbinding"
	_operationcontext "common/operationcontext"

	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestUnaryServerInterceptorBindsBeforeHandler(t *testing.T) {
	runtime := testRuntime(t, _operationbinding.DispositionDirect, nil)
	called := false

	_, err := runtime.UnaryServerInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{FullMethod: "/types.tqd.MapWorkspaceService/CreateGeneratedReport"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			called = true
			selection, ok := _operationcontext.FromContext(ctx)
			if !ok || selection.Code != testOperation {
				t.Fatalf("operation selection = %+v ok=%v", selection, ok)
			}
			return struct{}{}, nil
		},
	)
	if err != nil {
		t.Fatalf("UnaryServerInterceptor() error = %v", err)
	}
	if !called {
		t.Fatal("handler was not called")
	}
}

package operationgrpc

import (
	"context"
	"testing"

	"common/operation"
	_ "pb/types/tqd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryServerInterceptorBindsDeclaredOperationBeforeHandler(
	t *testing.T,
) {
	t.Parallel()

	called := false

	_, err := UnaryServerInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.MapWorkspaceService/CreateGeneratedReport",
		},
		func(ctx context.Context, req any) (any, error) {
			called = true

			code, ok := operation.FromContext(ctx)
			if !ok {
				t.Fatal("operation is missing from handler context")
			}
			if code != operation.Code("workspace.report.generate") {
				t.Fatalf(
					"operation = %q, want %q",
					code,
					"workspace.report.generate",
				)
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

func TestUnaryServerInterceptorLeavesUnannotatedMethodUnselected(
	t *testing.T,
) {
	t.Parallel()

	_, err := UnaryServerInterceptor(
		context.Background(),
		struct{}{},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.MapWorkspaceService/ListFollowedParcels",
		},
		func(ctx context.Context, req any) (any, error) {
			if code, ok := operation.FromContext(ctx); ok {
				t.Fatalf(
					"unexpected operation = %q",
					code,
				)
			}

			return struct{}{}, nil
		},
	)
	if err != nil {
		t.Fatalf("UnaryServerInterceptor() error = %v", err)
	}
}

func TestUnaryServerInterceptorKeepsSameOperationBinding(
	t *testing.T,
) {
	t.Parallel()

	ctx, err := operation.Bind(
		context.Background(),
		operation.Code("workspace.report.generate"),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = UnaryServerInterceptor(
		ctx,
		struct{}{},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.MapWorkspaceService/CreateGeneratedReport",
		},
		func(ctx context.Context, req any) (any, error) {
			code, ok := operation.FromContext(ctx)
			if !ok ||
				code != operation.Code("workspace.report.generate") {
				t.Fatalf(
					"operation = %q, %v",
					code,
					ok,
				)
			}

			return struct{}{}, nil
		},
	)
	if err != nil {
		t.Fatalf("UnaryServerInterceptor() error = %v", err)
	}
}

func TestUnaryServerInterceptorRejectsConflictingOperation(
	t *testing.T,
) {
	t.Parallel()

	ctx, err := operation.Bind(
		context.Background(),
		operation.Code("workspace.report.delete"),
	)
	if err != nil {
		t.Fatal(err)
	}

	called := false

	_, err = UnaryServerInterceptor(
		ctx,
		struct{}{},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.MapWorkspaceService/CreateGeneratedReport",
		},
		func(ctx context.Context, req any) (any, error) {
			called = true
			return struct{}{}, nil
		},
	)

	if status.Code(err) != codes.Internal {
		t.Fatalf(
			"status code = %v, want Internal; error = %v",
			status.Code(err),
			err,
		)
	}
	if called {
		t.Fatal("handler must not run after operation conflict")
	}
}

func TestUnaryServerInterceptorIgnoresRequestOperationField(
	t *testing.T,
) {
	t.Parallel()

	type request struct {
		Operation string
	}

	_, err := UnaryServerInterceptor(
		context.Background(),
		request{
			Operation: "workspace.report.delete",
		},
		&grpc.UnaryServerInfo{
			FullMethod: "/types.tqd.MapWorkspaceService/CreateGeneratedReport",
		},
		func(ctx context.Context, req any) (any, error) {
			code, ok := operation.FromContext(ctx)
			if !ok ||
				code != operation.Code("workspace.report.generate") {
				t.Fatalf(
					"operation = %q, %v",
					code,
					ok,
				)
			}

			return struct{}{}, nil
		},
	)
	if err != nil {
		t.Fatalf("UnaryServerInterceptor() error = %v", err)
	}
}

func TestUnaryServerInterceptorRejectsMissingMethodInfo(
	t *testing.T,
) {
	t.Parallel()

	called := false

	_, err := UnaryServerInterceptor(
		context.Background(),
		struct{}{},
		nil,
		func(ctx context.Context, req any) (any, error) {
			called = true
			return struct{}{}, nil
		},
	)

	if status.Code(err) != codes.Internal {
		t.Fatalf(
			"status code = %v, want Internal; error = %v",
			status.Code(err),
			err,
		)
	}
	if called {
		t.Fatal("handler must not run without method info")
	}
}

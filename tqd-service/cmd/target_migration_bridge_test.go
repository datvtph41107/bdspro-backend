package cmd

import (
	"common/businessoperation"
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestUsesTargetQuotaRequiresFeatureFlag(t *testing.T) {
	if usesTargetQuota(false, createGeneratedReportMethod) {
		t.Fatal("target quota must stay disabled when feature flag is off")
	}
	if !usesTargetQuota(true, createGeneratedReportMethod) {
		t.Fatal("target quota must handle migrated method when feature flag is on")
	}
	if usesTargetQuota(true, "/types.tqd.ParcelService/GetParcel") {
		t.Fatal("non-target RPC must keep legacy quota")
	}
}

func TestSkipTargetOperationRuntimeKeepsLegacyWhenFlagDisabled(t *testing.T) {
	legacyCalls := 0
	legacy := func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		legacyCalls++
		return handler(ctx, req)
	}
	handler := func(context.Context, interface{}) (interface{}, error) {
		return "ok", nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: createGeneratedReportMethod}

	if _, err := skipTargetOperationRuntime(legacy, false)(context.Background(), nil, info, handler); err != nil {
		t.Fatal(err)
	}
	if legacyCalls != 1 {
		t.Fatalf("legacy calls with flag off = %d, want 1", legacyCalls)
	}

	if _, err := skipTargetOperationRuntime(legacy, true)(context.Background(), nil, info, handler); err != nil {
		t.Fatal(err)
	}
	if legacyCalls != 1 {
		t.Fatalf("legacy calls with flag on = %d, want unchanged", legacyCalls)
	}
}

func TestBusinessOperationInterceptorFollowsFeatureFlag(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(map[bool]string{false: "disabled", true: "enabled"}[enabled], func(t *testing.T) {
			interceptor, err := newBusinessOperationInterceptor(enabled)
			if err != nil {
				t.Fatal(err)
			}

			_, err = interceptor(
				context.Background(),
				nil,
				&grpc.UnaryServerInfo{FullMethod: createGeneratedReportMethod},
				func(ctx context.Context, _ interface{}) (interface{}, error) {
					operation, ok := businessoperation.FromContext(ctx)
					if enabled && (!ok || operation != businessoperation.WorkspaceReportGenerate) {
						t.Fatalf("operation = %q, %v", operation, ok)
					}
					if !enabled && ok {
						t.Fatalf("operation must not be bound with flag off: %q", operation)
					}
					return nil, nil
				},
			)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

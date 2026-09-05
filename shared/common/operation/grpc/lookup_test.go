package operationgrpc

import (
	"errors"
	"testing"

	"common/operation"
	operationpb "pb/types/operation"
	_ "pb/types/tqd"

	"google.golang.org/protobuf/proto"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

func TestCodeForMethodReturnsDeclaredOperation(t *testing.T) {
	t.Parallel()

	code, found, err := codeForMethod(
		"/types.tqd.MapWorkspaceService/CreateGeneratedReport",
	)
	if err != nil {
		t.Fatalf("codeForMethod() error = %v", err)
	}
	if !found {
		t.Fatal("operation annotation was not found")
	}
	if code != operation.Code("workspace.report.generate") {
		t.Fatalf(
			"operation = %q, want %q",
			code,
			"workspace.report.generate",
		)
	}
}

func TestCodeForMethodLeavesUnannotatedMethodUnselected(t *testing.T) {
	t.Parallel()

	code, found, err := codeForMethod(
		"/types.tqd.MapWorkspaceService/ListFollowedParcels",
	)
	if err != nil {
		t.Fatalf("codeForMethod() error = %v", err)
	}
	if found {
		t.Fatalf("unexpected operation = %q", code)
	}
}

func TestCodeForMethodRejectsMalformedFullMethod(t *testing.T) {
	t.Parallel()

	for _, fullMethod := range []string{
		"",
		"types.tqd.MapWorkspaceService/CreateGeneratedReport",
		"/types.tqd.MapWorkspaceService",
		"/types.tqd.MapWorkspaceService/",
		"//CreateGeneratedReport",
		"/types.tqd.MapWorkspaceService/CreateGeneratedReport/extra",
	} {
		t.Run(fullMethod, func(t *testing.T) {
			t.Parallel()

			if _, _, err := codeForMethod(fullMethod); err == nil {
				t.Fatalf(
					"codeForMethod(%q) error = nil",
					fullMethod,
				)
			}
		})
	}
}

func TestCodeFromOptionsRejectsInvalidOperation(t *testing.T) {
	t.Parallel()

	options := &descriptorpb.MethodOptions{}
	proto.SetExtension(
		options,
		operationpb.E_Operation,
		&operationpb.Operation{
			Code: "Workspace.Report.Generate",
		},
	)

	_, found, err := codeFromOptions(options)
	if err == nil {
		t.Fatal("codeFromOptions() error = nil")
	}
	if found {
		t.Fatal("invalid operation must not be reported as found")
	}
	if !errors.Is(err, operation.ErrInvalidCode) {
		t.Fatalf(
			"codeFromOptions() error = %v, want ErrInvalidCode",
			err,
		)
	}
}

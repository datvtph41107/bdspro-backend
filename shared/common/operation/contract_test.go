package operation_test

import (
	"testing"

	"common/operation"
	operationpb "pb/types/operation"
	tqdpb "pb/types/tqd"

	"google.golang.org/protobuf/proto"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

func TestCreateGeneratedReportContractDeclaresOperation(t *testing.T) {
	t.Parallel()

	service := tqdpb.File_tqd_map_workspace_proto.
		Services().
		ByName("MapWorkspaceService")

	if service == nil {
		t.Fatal("MapWorkspaceService descriptor is missing")
	}

	method := service.
		Methods().
		ByName("CreateGeneratedReport")

	if method == nil {
		t.Fatal("CreateGeneratedReport descriptor is missing")
	}

	options, ok :=
		method.Options().(*descriptorpb.MethodOptions)

	if !ok || options == nil {
		t.Fatal(
			"CreateGeneratedReport method options are unavailable",
		)
	}

	if !proto.HasExtension(
		options,
		operationpb.E_Operation,
	) {
		t.Fatal(
			"CreateGeneratedReport operation annotation is missing",
		)
	}

	value := proto.GetExtension(
		options,
		operationpb.E_Operation,
	)

	annotation, ok :=
		value.(*operationpb.Operation)

	if !ok || annotation == nil {
		t.Fatalf(
			"operation annotation type = %T",
			value,
		)
	}

	code, err :=
		operation.Parse(annotation.GetCode())

	if err != nil {
		t.Fatalf(
			"operation.Parse(%q) error = %v",
			annotation.GetCode(),
			err,
		)
	}

	const want = "workspace.report.generate"

	if code != operation.Code(want) {
		t.Fatalf(
			"operation code = %q, want %q",
			code,
			want,
		)
	}
}

func TestUnannotatedWorkspaceMethodDoesNotDeclareOperation(
	t *testing.T,
) {
	t.Parallel()

	service := tqdpb.File_tqd_map_workspace_proto.
		Services().
		ByName("MapWorkspaceService")

	if service == nil {
		t.Fatal("MapWorkspaceService descriptor is missing")
	}

	method := service.
		Methods().
		ByName("ListFollowedParcels")

	if method == nil {
		t.Fatal(
			"ListFollowedParcels descriptor is missing",
		)
	}

	options, ok :=
		method.Options().(*descriptorpb.MethodOptions)

	if !ok || options == nil {
		t.Fatal(
			"ListFollowedParcels method options are unavailable",
		)
	}

	if proto.HasExtension(
		options,
		operationpb.E_Operation,
	) {
		t.Fatal(
			"ListFollowedParcels must remain unannotated in this slice",
		)
	}
}
